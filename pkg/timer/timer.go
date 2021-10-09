package timer

import (
	"github.com/robfig/cron"
	"github.com/zktnotify/zktnotify/pkg/config"
	"github.com/zktnotify/zktnotify/pkg/service"
)

func SetupTimer() {
	c := cron.New()
	c.AddFunc(config.Config.MonthDaily.Cron, service.SendMonthDaily)
	c.Start()
}
