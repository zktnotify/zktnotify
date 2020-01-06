package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/zktnotify/zktnotify/models"
	"github.com/zktnotify/zktnotify/pkg/config"
	"github.com/zktnotify/zktnotify/pkg/notify/typed"
	"github.com/zktnotify/zktnotify/pkg/zkt"
)

// Service main work service
func Service(ctx context.Context) {
	go func() {
		for {
			duration := cardDuration()

			select {
			case <-time.After(duration):
				if err := RetrieveCardTime(RetrieveAllUsers()); err != nil {
					log.Println(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		for {
			duration := sleepDuration()

			select {
			case <-time.After(time.Duration(duration) * time.Second):
				if err := CardTimeNotification(RetrieveWorkingUsers()); err != nil {
					log.Println(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func RetrieveAllUsers() []models.User {
	users, err := models.AllUsers()
	if err != nil {
		log.Println(err)
	}
	return users
}

func RetrieveWorkingUsers() []models.User {
	users, err := models.WorkingUsers()
	if err != nil {
		log.Println(err)
	}
	return users
}

func RetrieveCardTime(users []models.User) error {
	for _, user := range users {

		tag, err := getTodayCardTime(user)
		if err != nil {
			log.Println(err)
			continue
		}
		if tag == nil {
			continue
		}

		cardTimes, err := models.CardTimes(models.CardTime{
			UserID:   user.UserID,
			CardDate: time.Now().Format("2006-01-02"),
		})
		if err != nil {
			log.Println(err)
			continue
		}

		for ix, timeVal := range tag.CardTimes.EveryTime() {
			cardTime := models.CardTime{
				UserID:      tag.UserID,
				Times:       uint64(ix + 1),
				CardDate:    tag.CardDate,
				CardTime:    timeVal,
				BadgeNumber: tag.BadgeNumber,
			}

			if !cardTimeMatched(cardTimes, cardTime) {

				NewNotifier() <- NotifyMessage{
					UserID: tag.UserID,
					Name:   tag.Name,
					Date:   tag.CardDate,
					Time:   timeVal,
				}
				// FIXME: 发送报告是由协程去实现，如果执行过快，会导致发送多条
				time.Sleep(time.Millisecond * 300)

				if err := cardTime.Punched(); err != nil {
					log.Println(err)
				}
			}
		}
	}
	return nil
}

func getTodayCardTime(user models.User) (*models.TimeTag, error) {
	var err error

	if err := zkt.Login(config.Config.ZKTServer.URL.Login, user.JobID, user.Password); err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	if user.UserID == 0 {
		if user.UserID, err = zkt.GetUserID(config.Config.ZKTServer.URL.UserID); err != nil {
			return nil, fmt.Errorf("retrieve user id failed: %w", err)
		}
		user.UpdateUserID()
	}

	timeTag, err := zkt.GetTimeTag(config.Config.ZKTServer.URL.TimeTag, user.UserID, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("get time tag failed: %w", err)
	}
	if timeTag == nil {
		return nil, nil
	}

	tag := timeTag.Today()
	if tag == nil {
		return nil, nil
	}
	if tag.CardTimes.Len() < 1 {
		return nil, nil
	}
	return tag, nil
}

func cardTimeMatched(pattern []models.CardTime, match models.CardTime) bool {
	for _, card := range pattern {
		if card.UserID == match.UserID &&
			card.CardDate == match.CardDate &&
			card.CardTime == match.CardTime {
			return true
		}
	}
	return false
}

var lastNotifiedTime int64

func CardTimeNotification(users []models.User) error {
	ctype := uint64(Remind)
	cdate := time.Now().Format("2006-01-02")
	ctime := func() string { return time.Now().Format("15:04:05") }

	if !isWorkDate(cdate) {
		return nil
	}

	if ctime() < config.Config.WorkTime.End {
		return nil
	}

	if time.Now().Unix()-lastNotifiedTime < int64(config.Config.WorkEnd.NotificationTick) {
		return nil
	}

	for _, user := range users {
		if models.IsNotified(user.UserID, cdate, uint64(Worked)) {
			continue
		}
		if !models.CanNotify(user.UserID, cdate) {
			continue
		}

		dingtalk := ZKTNotifier{
			UID:        user.UserID,
			Name:       user.Name,
			Date:       cdate,
			Time:       ctime(),
			Type:       ctype,
			NotifyType: typed.NotifierType(user.NotifyType),
			Token:      user.NotifyToken,
		}
		dingtalk.send()
		models.UpdateNotice(user.UserID, ctype, cdate, ctime())
	}
	lastNotifiedTime = time.Now().Unix()
	return nil
}

func atou8(s string) uint8 {
	u, _ := strconv.ParseUint(s, 10, 64)
	return uint8(u)
}

func atou16(s string) uint16 {
	u, _ := strconv.ParseUint(s, 10, 64)
	return uint16(u)
}

func isWorkDate(cdate string) bool {
	cdates := strings.Split(cdate, "-")
	h, err := models.GetHoliday(atou16(cdates[0]), atou8(cdates[1]), atou8(cdates[2]))

	if err != nil {
		switch time.Now().Weekday() {
		case time.Sunday, time.Saturday:
			return false
		}
		return true
	}

	if !h.WorkDay {
		return false
	}
	return true
}

var firstNotified = true

func sleepDuration() int {
	defer func() { firstNotified = false }()

	var tick = 1

	nTime, err := time.ParseInLocation("2006-01-02 15:04:05",
		time.Now().Format("2006-01-02 ")+config.Config.WorkTime.End, time.Local)
	if err != nil {
		log.Println(err)
		return tick
	}

	dtime := nTime.Unix() - time.Now().Unix()
	if dtime >= 0 {
		if dtime < 5 {
			tick = 1
		} else {
			tick = int(dtime) - 5
		}
	} else {
		if firstNotified {
			tick = 1
		} else {
			tick = 30
		}
	}

	return tick
}

// cardDuration caculate access card time server interval
// using default value within 30 minutes before and after work or work end
// otherwise, access once in 10 minutes
func cardDuration() time.Duration {
	var (
		timeRange    = int64(30 * 60)
		defaultTick  = time.Duration(config.Config.TimeTick) * time.Second
		outRangeTick = 10 * time.Minute

		mktime = func(suffix string) (time.Time, error) {
			return time.ParseInLocation("2006/01/02 15:04:05",
				fmt.Sprintf("%s %s", time.Now().Format("2006/01/02"), suffix),
				time.Local)
		}
		inscope = func(wtime int64) bool {
			max, min := wtime, time.Now().Local().Unix()
			if max < min {
				max, min = min, max
			}
			return max-min < timeRange
		}
	)

	workTimeEnd, err1 := mktime(config.Config.WorkTime.End)
	workTimeStart, err2 := mktime(config.Config.WorkTime.Start)
	if err1 != nil || err2 != nil {
		err := err1
		if err == nil {
			err = err2
		}
		log.Println("mktime failed:", err)
		return defaultTick
	}

	if inscope(workTimeStart.Unix()) || inscope(workTimeEnd.Unix()) {
		return defaultTick
	}
	return outRangeTick
}
