package service

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/leaftree/onoffice/models"
	"github.com/leaftree/onoffice/pkg/config"
	"github.com/leaftree/onoffice/pkg/notify/dingtalk"
	"github.com/leaftree/onoffice/pkg/xerror"
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

type Notifier interface {
	Notify() error
	CanNotify() bool
}

type DingTalkNotifier struct {
	URL     string
	UID     uint64
	Name    string
	Date    string
	Time    string
	Type    uint64
	Account string
}

func (dtn *DingTalkNotifier) Notify() error {
	if dtn.CanNotify() {
		user := models.GetUser(dtn.UID)
		if user == nil {
			return fmt.Errorf("user(%d) not found", dtn.UID)
		}

		dtn.URL = user.NotifyURL
		dtn.Account = user.NotifyAccount

		models.Notified(dtn.UID, dtn.Type, dtn.Date, dtn.Time)
		return dtn.send()
	}
	return nil
}

func (dtn *DingTalkNotifier) CanNotify() bool {
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

func (dtn *DingTalkNotifier) send() error {
	return dingtalk.SendNotify(dtn.URL, dtn.msgTextTemplate(), dingtalk.Receiver{AtMobiles: []string{dtn.Account}})
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

		handler := &DingTalkNotifier{
			UID:  msg.UserID,
			Name: msg.Name,
			Date: msg.Date,
			Time: msg.Time,
			Type: uint64(ctype),
		}

		err := handler.Notify()
		xerror.LogError(err)
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

	card := cards[0]
	delay := uint32(0)
	for _, item := range config.Config.DelayWorkTime.Item {
		if card.CardTime < item.Time {
			break
		}
		delay = item.Delay
	}

	tt, _ := time.Parse("2006-01-02 15:04:05", cdate+" "+config.Config.WorkTime.Start)
	if ctime > tt.Add(time.Duration(delay)*time.Second).Format("15:04:05") {
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

func (dtn *DingTalkNotifier) msgTextTemplate() string {

	var msg string = "大兄弟，你已经打卡了，是上班、下班自己判断"
	templateText := map[uint64]string{
		Remind:    "{{.Name}}，该下班打卡了，当前时间{{.Date}} {{.Time}}",
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
