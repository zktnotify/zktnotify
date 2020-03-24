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
	xnotify "github.com/zktnotify/zktnotify/pkg/notify"
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

			punched := func() {
				if err := cardTime.Punched(); err != nil {
					log.Println(err)
				}
			}

			if !cardTimeMatched(cardTimes, cardTime) {
				NewNotifier() <- Notification{
					UserID:     tag.UserID,
					Name:       tag.Name,
					Date:       tag.CardDate,
					Time:       timeVal,
					AfterHooks: []HookFunc{punched},
				}
			}
		}
	}
	return nil
}

func getTodayCardTime(user models.User) (*models.TimeTag, error) {
	var err error

	if !zkt.HasCookie(user.JobID, user.UserID) {
		if err = zkt.Login(user.JobID, user.UserID, user.Password); err != nil {
			return nil, fmt.Errorf("%s login failed: %v", user.JobID, err)
		}

		if user.UserID == 0 {
			if user.UserID, err = zkt.GetUserID(user.JobID); err != nil {
				return nil, fmt.Errorf("retrieve user id failed: %w", err)
			}
			user.UpdateUserID()
		}
	}

	timeTag, err := zkt.GetTimeTag(user.UserID, time.Now(), time.Now())
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
	wtype := typed.OffWork
	status := typed.Remind
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
		// 特殊时期上班(2019-2020新冠)并且用户愿意接收下班打卡提醒
		if config.Config.IsSpecialPeriod && !user.SpecialPeriodNotify {
			continue
		}
		if !models.CanNotify(user.UserID, cdate) {
			continue
		}

		msg := typed.Message{
			UID:        user.UserID,
			Name:       user.Name,
			Date:       cdate,
			Time:       ctime(),
			Type:       wtype,
			Status:     status,
			Token:      user.NotifyToken,
			Account:    user.NotifyAccount,
			NotifyType: typed.NotifierType(user.NotifyType),
		}
		if err := sendNotice(msg); err != nil {
			log.Printf("send take card notify for user(%d) failed:%v", user.UserID, err)
			continue
		}
		models.UpdateNotice(user.UserID, uint64(status), cdate, ctime())
	}
	lastNotifiedTime = time.Now().Unix()
	return nil
}

func sendNotice(msg typed.Message) error {
	dtn := &Notification{
		UserID: msg.UID,
		Date:   msg.Date,
		Status: msg.Status,
	}

	sender := xnotify.New(msg)
	sender.SetCancelURL(dtn.shortURL())
	sender.SetAppToken(config.Config.XServer.NotificationServer.AppToken)

	receiver := typed.Receiver{
		All: false,
		ID:  []string{msg.Account},
	}
	return sender.Notify(msg.Token, sender.Template(msg), receiver)
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

	return h.WorkDay == 1
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

	if config.Config.Enviroment == "dev" {
		return defaultTick
	}

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
