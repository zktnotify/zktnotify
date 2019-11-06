package serverchan

import "testing"
import "github.com/leaftree/ctnotify/pkg/notify/typed"

func TestNotifier(t *testing.T) {
	var sender typed.Notifier
	sender = New()
	err := sender.Notify("https://sc.ftqq.com/SCU65820T6429cf7046f4c0d21f24e4ea3f254a025dc15b3d670c7.send", "刘云峰，你已经上班打卡，打卡时间2019-11-06 09:10:230")
	if err != nil {
		t.Log(err)
	}
}
