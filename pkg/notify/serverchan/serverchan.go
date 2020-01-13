package serverchan

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	httpurl "net/url"
	"text/template"

	"github.com/zktnotify/zktnotify/pkg/notify/typed"
)

const (
	NotifyHost = "https://sc.ftqq.com"
)

type ServerChan struct{}

func New() typed.Notifier {
	return &ServerChan{}
}

type responsed struct {
	ErrNo   int    `json:"errno"`
	ErrMsg  string `json:"errmsg"`
	DataSet string `json:"dataset"`
}

func (s *ServerChan) Notify(token string, msg string, receiver ...typed.Receiver) error {
	if s == nil || token == "" {
		return nil
	}

	url := URL(token)
	payload := httpurl.Values{}
	payload.Set("text", msg)
	payload.Set("desp", "# 测试啦")

	rsp, err := http.Get(url + "?" + payload.Encode())
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	response := responsed{}
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}
	if response.ErrNo != 0 {
		return errors.New(response.ErrMsg)
	}

	return nil
}

func (s *ServerChan) SetCancelURL(url string)  {}
func (s *ServerChan) SetAppToken(token string) {}

func (s *ServerChan) Template(msg typed.Message) string {
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
	return NotifyHost + "/" + token + ".send"
}
