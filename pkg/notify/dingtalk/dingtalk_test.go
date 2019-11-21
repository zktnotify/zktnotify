package dingtalk

import "testing"
import "github.com/zktnotify/zktnotify/pkg/notify/typed"

func TestNotifier(t *testing.T) {
	var sender typed.Notifier
	sender = New()
	sender.Notify("4e35556a0ef4fb9fdba399e147df9a533b4bb19f29919dd88e906ababc35f5c3",
		"测试", typed.Receiver{All: false, ID: []string{"15920385660"}})
}
