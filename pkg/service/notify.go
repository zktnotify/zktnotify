package service

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"sort"
	"text/template"
	"time"

	"github.com/zktnotify/zktnotify/models"
	"github.com/zktnotify/zktnotify/pkg/config"
	xnotify "github.com/zktnotify/zktnotify/pkg/notify"
	"github.com/zktnotify/zktnotify/pkg/notify/typed"
	"github.com/zktnotify/zktnotify/pkg/shorturl"
)

const (
	Invalid = iota + 0
	ToWork
	Worked
	Midway
	Lated
	Remind    // notify to take a card
	DelayWork // delay in work
)

type NotifyMessage struct {
	UserID uint64
	Name   string
	Date   string
	Time   string
}

type ZKTNotifier struct {
	UID        uint64
	Name       string
	Date       string
	Time       string
	Type       uint64
	Token      string
	account    string
	NotifyType typed.NotifierType
}

func (dtn *ZKTNotifier) Notify() error {
	if dtn.CanNotify() {
		user := models.GetUser(dtn.UID)
		if user == nil {
			return fmt.Errorf("user(%d) not found", dtn.UID)
		}

		dtn.Token = user.NotifyToken
		dtn.account = user.NotifyAccount
		dtn.NotifyType = typed.NotifierType(user.NotifyType)

		err := models.Notified(dtn.UID, dtn.Type, dtn.Date, dtn.Time)
		if err != nil {
			return err
		}
		return dtn.send()
	}
	return nil
}

func (dtn *ZKTNotifier) CanNotify() bool {
	switch dtn.Type {
	case Invalid, ToWork, Worked, Lated, Remind, DelayWork:
	default:
		return false
	}

	if models.IsNotified(dtn.UID, dtn.Date, dtn.Type) {
		return false
	}
	return true
}

func (dtn *ZKTNotifier) send() error {
	if !typed.Valid(dtn.NotifyType) {
		return errors.New("invalid notification service type")
	}
	return xnotify.New(dtn.NotifyType).Notify(dtn.Token, dtn.msgTextTemplate(), typed.Receiver{All: false, ID: []string{dtn.account}})
}

func NewNotifier() chan<- NotifyMessage {
	ch := make(chan NotifyMessage)

	go func() {
		msg := <-ch
		early, last := models.EarliestAndLatestCardTime(msg.UserID, msg.Date)
		ctype := cardTimeType(early, last, msg.Time)

		if ctype == Lated && delayInWork(msg.UserID, msg.Date, msg.Time) {
			ctype = DelayWork
		}

		handler := &ZKTNotifier{
			UID:  msg.UserID,
			Name: msg.Name,
			Date: msg.Date,
			Time: msg.Time,
			Type: uint64(ctype),
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

func cardTimeType(early, last *models.CardTime, ctime string) uint64 {
	if early == nil || last == nil { // card time not found
		if ctime > config.Config.WorkTime.End { // Work end
			return Invalid
		}
		if ctime > config.Config.WorkTime.Start { // first card and after starting-work time, you late
			return Lated
		}
		return ToWork // normal card
	}

	if ctime <= early.CardTime { // ToWork
		return ToWork
	}

	if ctime < config.Config.WorkTime.End { // working
		return Midway
	}
	if ctime > config.Config.WorkTime.End { // work end
		return Worked
	}

	return Midway
}

func (dtn *ZKTNotifier) shortURL() string {
	return shorturl.ShortURL("/counternotice", map[string]interface{}{
		"userid":    dtn.UID,
		"card_date": dtn.Date,
	})
}

func (dtn *ZKTNotifier) msgTextTemplate() string {

	var msg string = "大兄弟，你已经打卡了，是上班、下班自己判断"
	templateText := map[uint64]string{
		Remind:    "{{.Name}}，该下班打卡了，当前时间{{.Date}} {{.Time}} " + dtn.shortURL(),
		ToWork:    "{{.Name}}，你已经上班打卡，打卡时间{{.Date}} {{.Time}}",
		Worked:    "{{.Name}}，你已经下班打卡，打卡时间{{.Date}} {{.Time}}",
		Lated:     "{{.Name}}，你已经上班打卡，打卡时间{{.Date}} {{.Time}}，可惜你迟到了",
		Invalid:   "{{.Name}}，你已经打卡，打卡时间{{.Date}} {{.Time}}，可是这个时候你打卡干嘛呢",
		DelayWork: "{{.Name}}，你已经上班打卡，打卡时间{{.Date}} {{.Time}}，昨晚下班有点晚，今天不迟到",
	}

	temp, ok := templateText[dtn.Type]
	if !ok {
		return msg
	}

	t, err := template.New("fylos").Parse(temp)
	if err != nil {
		return msg
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, dtn); err != nil {
		return msg
	}
	return buf.String()
}
