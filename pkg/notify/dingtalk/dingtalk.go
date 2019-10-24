package dingtalk

import (
	"bytes"
	"encoding/json"
	"net/http"
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

func SendNotify(url string, msg string, at ...Receiver) error {

	if url == "" {
		return nil
	}

	Dmsg := Message{
		Type: "text",
		Text: Text{Content: msg},
	}
	if len(at) == 1 {
		Dmsg.At = at[0]
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
	defer resp.Body.Close()

	return nil
}
