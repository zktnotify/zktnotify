package dingtalk

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/alecthomas/template"
	"github.com/zktnotify/zktnotify/pkg/notify/typed"
)

const (
	NotifyHost = "https://oapi.dingtalk.com/robot/send"
)

type Text struct {
	Content string `json:"content"`
}

type Receiver struct {
	IsAtAll   bool     `json:"isAtAll"`
	AtMobiles []string `json:"atMobiles"`
}

type Message struct {
	Type string   `json:"msgtype"`
	Text Text     `json:"text"`
	At   Receiver `json:"at"`
}

type DingTalk struct {
}

func New() typed.Notifier {
	return &DingTalk{}
}

func (d *DingTalk) Notify(token string, msg string, receiver ...typed.Receiver) error {

	if d == nil || token == "" {
		return nil
	}

	Dmsg := Message{
		Type: "text",
		Text: Text{Content: msg},
	}
	url := URL(token)

	if len(receiver) == 1 {
		who := Receiver{}
		who.IsAtAll = receiver[0].All
		who.AtMobiles = append(who.AtMobiles, receiver[0].ID...)
		Dmsg.At = who
	}
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(Dmsg)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

func (d *DingTalk) SetCancelURL(url string)  {}
func (d *DingTalk) SetAppToken(token string) {}

func (d *DingTalk) Template(msg typed.Message) string {

	var text string = "大兄弟，你已经打卡了，是上班、下班自己判断"
	templateText := map[typed.TemplateID]string{
		typed.Remind:    "{{.Name}}，该下班打卡了，当前时间{{.Date}} {{.Time}} " + msg.CancelURL,
		typed.ToWork:    "{{.Name}}，你已经上班打卡，打卡时间{{.Date}} {{.Time}}",
		typed.Worked:    "{{.Name}}，你已经下班打卡，打卡时间{{.Date}} {{.Time}}",
		typed.Lated:     "{{.Name}}，你已经上班打卡，打卡时间{{.Date}} {{.Time}}，可惜你迟到了",
		typed.Invalid:   "{{.Name}}，你已经打卡，打卡时间{{.Date}} {{.Time}}，可是这个时候你打卡干嘛呢",
		typed.DelayWork: "{{.Name}}，你已经上班打卡，打卡时间{{.Date}} {{.Time}}，昨晚下班有点晚，今天不迟到",
	}

	temp, ok := templateText[msg.Status]
	if !ok {
		return text
	}

	t, err := template.New("fylos").Parse(temp)
	if err != nil {
		return text
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, msg); err != nil {
		return text
	}
	return buf.String()
}

func URL(token string) string {
	return NotifyHost + "?access_token=" + token
}
