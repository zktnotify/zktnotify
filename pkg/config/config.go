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
	UserID   int    `json:"uid"`
	UserName string `json:"user"`
	Password string `json:"password"`
	TimeTick uint32 `json:"interval"`
	WorkTime struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"worktime"`
	ZKTServer struct {
		URL struct {
			Login   string `json:"login"`
			TimeTag string `json:"timeTag"`
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

	if err := load(filename); err != nil {
		return config{}, err
	}
	if err := Config.Validator(); err != nil {
		return config{}, err
	}
	return Config, nil
}

func load(filename string) error {
	Config = config{
		UserID:   3494,
		UserName: "1905",
		Password: "123kbc,./",
		TimeTick: 2,
	}
	Config.WorkTime.End = "18:00:00"
	Config.WorkTime.Start = "09:15:59"
	Config.ZKTServer.URL.Login = "http://money.fylos.cn:1234/selfservice/login/" // XXX:host route-path split
	Config.ZKTServer.URL.TimeTag = "http://money.fylos.cn:1234/grid/att/CardTimes/"
	Config.XServer.Addr = "http://127.0.0.1:4444"
	Config.XServer.File.Pid = filepath.Join(WorkDir, AppName+".pid")
	Config.XServer.File.Log = filepath.Join(WorkDir, AppName+".log")
	Config.XServer.File.DB = filepath.Join(WorkDir, "data", "data.db")

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		data = []byte(`{}`)
	}
	if err := json.Unmarshal(data, &Config); err != nil {
		return err
	}

	dir := filepath.Dir(filename)
	if f, err := os.Stat(dir); !(err == nil && f.IsDir()) {
		os.MkdirAll(dir, 0755)
	}
	data, _ = json.MarshalIndent(Config, "", "\t")
	ioutil.WriteFile(filename, data, 0644)

	return nil
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
