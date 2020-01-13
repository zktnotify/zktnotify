package notify

import (
	"github.com/zktnotify/zktnotify/pkg/notify/dingtalk"
	"github.com/zktnotify/zktnotify/pkg/notify/serverchan"
	"github.com/zktnotify/zktnotify/pkg/notify/typed"
	"github.com/zktnotify/zktnotify/pkg/notify/wxpusher"
)

func New(msg typed.Message) typed.Notifier {
	ntype := msg.NotifyType

	switch ntype {
	case typed.DINGTALK:
		return dingtalk.New()
	case typed.SERVERCHAN:
		return serverchan.New()
	case typed.WXPUSHER:
		return wxpusher.New()
	}
	return nil
}
