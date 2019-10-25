package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/leaftree/onoffice/pkg/xpath"
)

var (
	AppName = "pin"
	WorkDir = filepath.Join(xpath.HomeDir(), "."+AppName)
	Config  config
)

type config struct {
	TimeTick uint32 `json:"interval"`
	WorkEnd  struct {
		NotificationTime string `json:"time"`
		NotificationTick uint32 `json:"interval"`
	} `json:"workend"`
	WorkTime struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"worktime"`
	ZKTServer struct {
		URL struct {
			Login   string `json:"login"`
			UserID  string `json:"userid"`
			TimeTag string `json:"timetag"`
		} `json:"url"`
	} `json:"zktserver"`
	XServer struct {
		Addr string `json:"addr"`
		Name string `json:"name"`
		File struct {
			Pid string `json:"pid"`
			Log string `json:"log"`
			DB  string `json:"database"`
		} `json:"file"`
	} `json:"xserver"`
}

func NewConfig(file ...string) (config, error) {
	filename := filepath.Join(WorkDir, "config.json")
	if len(file) != 0 && file[0] != "" {
		filename = file[0]
	}

	cfg, err := load(filename)
	if err != nil {
		return config{}, err
	}
	if err := Config.Validator(); err != nil {
		return config{}, err
	}
	Config = *cfg
	return Config, nil
}

func load(filename string) (*config, error) {
	cfg := &config{
		TimeTick: 2,
	}
	cfg.WorkTime.End = "18:00:00"
	cfg.WorkTime.Start = "09:15:59"
	cfg.WorkEnd.NotificationTick = 1800
	cfg.WorkEnd.NotificationTime = "18:00:00"
	cfg.ZKTServer.URL.UserID = "http://money.fylos.cn:1234/selfservice/selfreport/"
	cfg.ZKTServer.URL.Login = "http://money.fylos.cn:1234/selfservice/login/" // XXX:host route-path split
	cfg.ZKTServer.URL.TimeTag = "http://money.fylos.cn:1234/grid/att/CardTimes/"
	cfg.XServer.Addr = "http://127.0.0.1:4444"
	cfg.XServer.File.Pid = filepath.Join(WorkDir, AppName+".pid")
	cfg.XServer.File.Log = filepath.Join(WorkDir, AppName+".log")
	cfg.XServer.File.DB = filepath.Join(WorkDir, "data", "data.db")

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		data = []byte(`{}`)
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	dir := filepath.Dir(filename)
	if f, err := os.Stat(dir); !(err == nil && f.IsDir()) {
		os.MkdirAll(dir, 0755)
	}
	data, _ = json.MarshalIndent(cfg, "", "\t")
	ioutil.WriteFile(filename, data, 0644)

	return cfg, nil
}

func (cfg config) Validator() error {
	if cfg.TimeTick < 1 && cfg.TimeTick > 30*60 {
		return errors.New("config.interval should be from 1 to 1800")
	}
	return nil
}

func (cfg config) LogName() string {
	return cfg.XServer.File.Log
}
