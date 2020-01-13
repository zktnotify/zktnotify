package service

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/zktnotify/zktnotify/models"
	"github.com/zktnotify/zktnotify/pkg/config"
	xnotify "github.com/zktnotify/zktnotify/pkg/notify"
	"github.com/zktnotify/zktnotify/pkg/notify/typed"
	"github.com/zktnotify/zktnotify/pkg/shorturl"
)

type CardStatus = typed.TemplateID
type CardType = typed.WorkType

type Notification struct {
	UserID     uint64
	Name       string
	Date       string
	Time       string
	Token      string
	Account    string
	Type       CardType
	Status     CardStatus
	NotifyType typed.NotifierType
}

func (dtn *Notification) Notify() error {
	if !dtn.CanNotify() {
		return nil
	}

	user := models.GetUser(dtn.UserID)
	if user == nil {
		return fmt.Errorf("user(%d) not found", dtn.UserID)
	}

	if err := models.Notified(dtn.UserID, uint64(dtn.Status), dtn.Date, dtn.Time); err != nil {
		return err
	}

	if !typed.Valid(dtn.NotifyType) {
		return errors.New("invalid notification service type")
	}

	msg := typed.Message{
		UID:        dtn.UserID,
		Date:       dtn.Date,
		Time:       dtn.Time,
		Type:       dtn.Type,
		Status:     dtn.Status,
		Name:       user.Name,
		Token:      user.NotifyToken,
		Account:    user.NotifyAccount,
		NotifyType: typed.NotifierType(user.NotifyType),
	}
	sender := xnotify.New(msg)

	sender.SetCancelURL(dtn.shortURL())
	sender.SetAppToken(config.Config.XServer.NotificationServer.AppToken)

	receiver := typed.Receiver{
		All: false,
		ID:  []string{dtn.Account},
	}
	return sender.Notify(msg.Token, sender.Template(msg), receiver)
}

func (dtn *Notification) CanNotify() bool {
	switch dtn.Status {
	case typed.Invalid, typed.ToWork, typed.Worked, typed.Lated, typed.Remind, typed.DelayWork:
	default:
		return false
	}

	if models.IsNotified(dtn.UserID, dtn.Date, uint64(dtn.Status)) {
		return false
	}
	return true
}

func (dtn *Notification) send() error {
	/*
		sender := xnotify.New(dtn.NotifyType)
		sender.SetAppToken(config.Config.XServer.NotificationServer.AppToken)

		if dtn.Status == typed.Remind {
			sender.SetCancelURL(dtn.shortURL())
		}

		return sender.Notify(
			dtn.Token,
			xnotify.Template(sender),
			typed.Receiver{
				All: false,
				ID:  []string{dtn.Account},
			},
		)
	*/
	return nil
}

func NewNotifier() chan<- Notification {
	ch := make(chan Notification)

	go func() {
		msg := <-ch
		early, last := models.EarliestAndLatestCardTime(msg.UserID, msg.Date)
		status := cardTimeStatus(early, last, msg.Time)
		wtype := workType(early, last)

		if status == typed.Lated && delayInWork(msg.UserID, msg.Date, msg.Time) {
			status = typed.DelayWork
		}

		handler := &Notification{
			UserID: msg.UserID,
			Name:   msg.Name,
			Date:   msg.Date,
			Time:   msg.Time,
			Type:   wtype,
			Status: status,
		}

		if err := handler.Notify(); err != nil {
			log.Println(err)
		}
	}()
	return ch
}

func delayInWork(uid uint64, cdate, ctime string) bool {
	t, _ := time.Parse("2006-01-02 15:04:05", cdate+" "+ctime)
	cards, err := models.CardTimes(models.CardTime{
		UserID:   uid,
		CardDate: time.Unix(t.Unix()-24*60*60, 0).Format("2006-01-02"),
	})

	if len(cards) == 0 || err != nil {
		return false
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i].CardTime > cards[j].CardTime })

	var (
		card  = cards[0]
		delay = uint32(0)
		items = make([]struct {
			Time  string `json:"time"`
			Delay uint32 `json:"delay"`
		}, len(config.Config.DelayWorkTime.Item))
	)
	copy(items, config.Config.DelayWorkTime.Item)
	sort.Slice(items, func(i, j int) bool { return items[i].Time < items[j].Time })

	for _, item := range items {
		if card.CardTime < item.Time {
			break
		}
		delay = item.Delay
	}

	tt, _ := time.Parse("2006-01-02 15:04:05", cdate+" "+config.Config.WorkTime.Start)
	if ctime > tt.Add(time.Duration(delay)*time.Minute).Format("15:04:05") {
		return false
	}
	return true
}

func workType(early, last *models.CardTime) CardType {
	if early == nil {
		return typed.Working
	}
	return typed.OffWork
}

func cardTimeStatus(early, last *models.CardTime, ctime string) CardStatus {
	if early == nil || last == nil { // card time not found
		if ctime > config.Config.WorkTime.End { // Work end
			return typed.Invalid
		}
		if ctime > config.Config.WorkTime.Start { // first card and after starting-work time, you late
			return typed.Lated
		}
		return typed.ToWork // normal card
	}

	if ctime <= early.CardTime { // ToWork
		return typed.ToWork
	}

	if ctime < config.Config.WorkTime.End { // working
		return typed.Midway
	}
	if ctime > config.Config.WorkTime.End { // work end
		return typed.Worked
	}

	return typed.Midway
}

func (dtn *Notification) shortURL() string {
	if dtn.Status != typed.Remind {
		return ""
	}
	return shorturl.ShortURL("/counternotice", map[string]interface{}{
		"userid":    dtn.UserID,
		"card_date": dtn.Date,
	})
}
