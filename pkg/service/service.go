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
}

func RetrieveAllUsers() []models.User {
	users, err := models.AllUsers()
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

		for _, timeVal := range tag.CardTimes.EveryTime() {
			cardTime := models.CardTime{
				UserID:      tag.UserID,
				Times:       uint64(tag.Times),
				CardDate:    tag.CardDate,
				CardTime:    timeVal,
				BadgeNumber: tag.BadgeNumber,
			}

			NewNotifier() <- NotifyMessage{
				UserID: tag.UserID,
				Date:   tag.CardDate,
				Time:   timeVal,
			}

			if !cardTimeMatched(cardTimes, cardTime) {
				if err := cardTime.Punched(); err != nil {
					log.Println(err)
				}
			}
		}
	}
	return nil
}

func getTodayCardTime(user models.User) (*models.TimeTag, error) {
	if err := zkt.Login(config.Config.ZKTServer.URL.Login, user.JobID, user.Password); err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
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
