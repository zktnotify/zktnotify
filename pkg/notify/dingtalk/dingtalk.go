package dingtalk

import (
	"bytes"
	"encoding/json"
	"net/http"

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

type DingTalk struct{}

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

func URL(token string) string {
	return NotifyHost + "?access_token=" + token
}
