package wxpusher

import (
	"testing"

	"github.com/zktnotify/zktnotify/pkg/notify/typed"
)

func TestMain(t *testing.T) {
	notifier := New()
	notifier.SetAppToken("AT_4EhDwiLfZfre2eYGWRfiPoeFkNlbciIW")

	msg := typed.Message{
		Date:   "2020/01/13",
		Time:   "10:30:00",
		Type:   typed.Working,
		Status: typed.ToWork,
	}
	t.Log(notifier.Notify("UID_XK7Qp5fAPTxNgszplAEAqgnOTebX", notifier.Template(msg)))

	msg = typed.Message{
		Date:   "2020/01/13",
		Time:   "10:30:00",
		Type:   typed.Working,
		Status: typed.Lated,
	}
	t.Log(notifier.Notify("UID_XK7Qp5fAPTxNgszplAEAqgnOTebX", notifier.Template(msg)))

}
