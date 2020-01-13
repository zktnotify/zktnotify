package wxpusher

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/zktnotify/zktnotify/pkg/notify/typed"
)

const (
	NotifyHost = "http://wxpusher.zjiecode.com"
	NotifyHook = NotifyHost + "/api/send/message"
)

// contentType 消息内容类型
type contentType int

const (
	ContentText = iota + 1
	ContentHtml
	ContentMarkdown
)

var (
	errlist = map[int]string{
		1000: "处理成功",
		1001: "业务异常错误",
		1002: "未认证",
		1003: "签名错误",
		1004: "接口不存在",
		1005: "服务器内部错误",
		1006: "和微信交互的过程中发生异常",
		1007: "网络异常",
		1008: "数据异常",
		1009: "未知异常",
		9999: "未知异常",
	}
)

func Error(code int) string {
	if msg, ok := errlist[code]; ok {
		return msg
	}
	return errlist[9999]
}

type responsed struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		UID       string `json:"uid"`
		TopicID   string `json:"topicId"`
		MessageID int    `json:"messageId"`
		Code      int    `json:"code"`
		Status    string `json:"status"`
	} `json:"data"`
	Success bool `json:"success"`
}

func (r responsed) Ok() bool {
	return r.Code == 1000 && r.Success
}

func (r responsed) Error() string {
	m1, m2 := Error(r.Code), ""
	if len(r.Data) > 0 {
		m2 = ": " + Error(r.Data[0].Code)
	}
	return m1 + m2
}

type WXPusher struct {
	AppToken    string      `json:"appToken"`
	Content     string      `json:"content"`
	ContentType contentType `json:"contentType"`
	TopicIDs    []int       `json:"topicIds"`
	UIDs        []string    `json:"uids"`
	URL         string      `json:"url"`
}

var _ typed.Notifier = (*WXPusher)(nil)

func New() typed.Notifier {
	notifier := WXPusher{
		ContentType: ContentMarkdown,
		TopicIDs:    []int{},
	}
	return &notifier
}

func (w *WXPusher) Notify(userToken string, msg string, receiver ...typed.Receiver) error {
	w.Content = msg
	w.UIDs = []string{userToken}

	data, err := json.Marshal(w)
	if err != nil {
		return err
	}

	resp, err := http.Post(NotifyHook, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result responsed
	err = json.Unmarshal(rData, &result)
	if err != nil {
		return err
	}
	if result.Ok() {
		return nil
	}
	return errors.New(result.Error())
}

func (w *WXPusher) SetCancelURL(url string) {
	w.URL = url
}

func (w *WXPusher) SetAppToken(token string) {
	w.AppToken = token
}

func (w *WXPusher) Template(msg typed.Message) string {
	tpl := struct {
		CardType  string
		Status    string
		Date      string
		Time      string
		CancelURL string
	}{
		CardType:  msg.Type.String(),
		Status:    msg.Status.String(),
		Date:      msg.Date,
		Time:      msg.Time,
		CancelURL: msg.CancelURL,
	}

	var defaultText = "oops ...."

	var normalText = `
# 打卡通知 ({{.CardType}})
* 状态：{{.Status}}
* 时间：{{.Time}}
* 日期：{{.Date}}
`

	var lateText = `
# 打卡通知 ({{.CardType}})
* 状态：{{.Status}}
* 时间：{{.Time}}
* 日期：{{.Date}}

![](http://pic.17qq.com/img_biaoqing/20262164.jpeg)
`

	var delayText = `
# 打卡通知 ({{.CardType}})
* 状态：{{.Status}}
* 时间：{{.Time}}
* 日期：{{.Date}}

【温馨提示】：昨晚辛苦了，今天没有迟到
`

	var remindText = `
# 该下班了，记得打卡
* 时间：{{.Time}}
* 日期：{{.Date}}

![](http://pic.17qq.com/img_biaoqing/20262164.jpeg)[点我取消]({{.CancelURL}})
`

	var text string
	switch msg.Status {
	case typed.Remind:
		text = remindText
	case typed.DelayWork:
		text = delayText
	case typed.Lated:
		text = lateText
	default:
		text = normalText
	}

	t, err := template.New("fylos").Parse(text)
	if err != nil {
		return defaultText
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, tpl); err != nil {
		return defaultText
	}

	return buf.String()
}
