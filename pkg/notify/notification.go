package notify

import (
	"github.com/leaftree/ctnotify/pkg/notify/dingtalk"
	"github.com/leaftree/ctnotify/pkg/notify/serverchan"
	"github.com/leaftree/ctnotify/pkg/notify/typed"
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
