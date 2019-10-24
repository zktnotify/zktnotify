package service

import (
	"bytes"
	"text/template"

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
	Remind // notify to take a card
)

type NotifyMessage struct {
	UserID uint64
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
		models.Notified(dtn.UID, dtn.Type, dtn.Date, dtn.Time)
		return dtn.send()
	}
	return nil
}

func (dtn *DingTalkNotifier) CanNotify() bool {
	user := models.GetUser(dtn.UID)
	if user == nil {
		return false
	}

	switch dtn.Type {
	case Invalid, ToWork, Worked, Lated:
	default:
		return false
	}

	if models.IsNotified(dtn.UID, dtn.Date, dtn.Type) {
		return false
	}

	dtn.URL = user.NotifyURL
	dtn.Name = user.Name
	dtn.Account = user.NotifyAccount
	return true
}

func (dtn *DingTalkNotifier) send() error {
	return dingtalk.SendNotify(dtn.URL, dtn.msgTextTemplate(), dingtalk.Receiver{AtMobiles: []string{dtn.Account}})
}

func NewNotifier() chan<- NotifyMessage {
	ch := make(chan NotifyMessage, 10)

	go func() {
		msg := <-ch
		early, last := models.EarliestAndLatestCardTime(msg.UserID, msg.Date)
		ctype := cardTimeType(early, last, msg.Time)

		handler := &DingTalkNotifier{
			UID:  msg.UserID,
			Date: msg.Date,
			Time: msg.Time,
			Type: uint64(ctype),
		}

		err := handler.Notify()
		xerror.LogError(err)
	}()
	return ch
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
		ToWork:  "{{.Name}}，你已经上班打卡，打卡时间{{.Date}} {{.Time}}",
		Worked:  "{{.Name}}，你已经下班打卡，打卡时间{{.Date}} {{.Time}}",
		Lated:   "{{.Name}}，你已经上班打卡，打卡时间{{.Date}} {{.Time}}，可惜你迟到了",
		Invalid: "{{.Name}}，你已经打卡，打卡时间{{.Date}} {{.Time}}，可是这个时候你打卡干嘛呢",
	}

	t, err := template.New("fylos").Parse(templateText[dtn.Type])
	if err != nil {
		return msg
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, dtn); err != nil {
		return msg
	}
	return buf.String()
}
