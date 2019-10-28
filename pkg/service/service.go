package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/leaftree/onoffice/models"
	"github.com/leaftree/onoffice/pkg/config"
	"github.com/leaftree/onoffice/pkg/zkt"
)

// Service main work service
func Service(ctx context.Context) {
	go func() {
		for {
			select {
			case <-time.After(time.Duration(config.Config.TimeTick) * time.Second):
				if err := RetrieveCardTime(RetrieveAllUsers()); err != nil {
					log.Println(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		var tick = 1

		nTime, err := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 ")+config.Config.WorkEnd.NotificationTime, time.Local)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}

		dtime := nTime.Unix() - time.Now().Unix()
		if dtime < -10 {
			tick = 60
		} else if dtime <= 5 {
			tick = 1
		} else if dtime > 60 {
			tick = 60
		} else {
			tick = 5
		}

		select {
		case <-time.After(time.Duration(tick) * time.Second):
			if err := CardTimeNotification(RetrieveWorkingUsers()); err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}
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
				if err := cardTime.Punched(); err != nil {
					log.Println(err)
				}

				NewNotifier() <- NotifyMessage{
					UserID: tag.UserID,
					Name:   tag.Name,
					Date:   tag.CardDate,
					Time:   timeVal,
				}
				// FIXME: 发送报告是由协程去实现，如果执行过快，会导致发送多条
				time.Sleep(time.Millisecond * 300)
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

func CardTimeNotification(users []models.User) error {
	ctype := uint64(Remind)
	cdate := time.Now().Format("2006-01-02")
	ctime := func() string { return time.Now().Format("15:04:05") }

	switch time.Now().Weekday() {
	case time.Sunday, time.Saturday:
		return nil
	}

	if ctime() < config.Config.WorkEnd.NotificationTime {
		fmt.Println("时间未到")
		return nil
	}

	for _, user := range users {
		// checking workend card time
		if models.IsNotified(user.UserID, cdate, uint64(Worked)) {
			continue
		}

		dingtalk := DingTalkNotifier{
			URL:     user.NotifyURL,
			UID:     user.UserID,
			Name:    user.Name,
			Date:    cdate,
			Time:    ctime(),
			Type:    ctype,
			Account: user.NotifyAccount,
		}
		fmt.Println("send")
		fmt.Println(dingtalk.send())
	}
	return nil
}