package notify

import (
	"github.com/zktnotify/zktnotify/pkg/notify/dingtalk"
	"github.com/zktnotify/zktnotify/pkg/notify/serverchan"
	"github.com/zktnotify/zktnotify/pkg/notify/typed"
)

func New(ntype typed.NotifierType) typed.Notifier {
	switch ntype {
	case typed.DINGTALK:
		return dingtalk.New()
	case typed.SERVERCHAN:
		return serverchan.New()
	}
	return nil
}
