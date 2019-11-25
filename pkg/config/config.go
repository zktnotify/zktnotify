package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/zktnotify/zktnotify/pkg/xpath"
)

var (
	AppName = "pin"
	WorkDir = filepath.Join(xpath.HomeDir(), "."+AppName)
	Config  config
)

type config struct {
	TimeTick uint32 `json:"interval"`
	WorkEnd  struct {
		NotificationTick uint32 `json:"interval"`
	} `json:"workend"`
	WorkTime struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"worktime"`
	DelayWorkTime struct {
		Item []struct {
			Time  string `json:"time"`
			Delay uint32 `json:"delay"`
		} `json:"item"`
	} `json:"delayworktime"`
	ZKTServer struct {
		URL struct {
			Login   string `json:"login"`
			UserID  string `json:"userid"`
			TimeTag string `json:"timetag"`
		} `json:"url"`
	} `json:"zktserver"`
	XClient struct {
		Server struct {
			Addr string `json:"addr"`
		} `json:"server"`
	} `json:"xclient"`
	XServer struct {
		Addr string `json:"addr"`
		Name string `json:"name"`
		File struct {
			Pid string `json:"pid"`
			Log string `json:"log"`
		} `json:"file"`
		DB struct {
			Type     string `json:"type"`
			Host     string `json:"host"`
			User     string `json:"user"`
			Name     string `json:"db_name"`
			Password string `json:"password"`
			Path     string `json:"path"`
		} `json:"database"`
		ShortURL struct {
			Server struct {
				AppKey  string `json:"appkey"`
				ApiAddr string `json:"address"`
			} `json:"server"`
			PrefixURL string `json:"prefixurl"`
		} `shorturl`
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
	if err := cfg.Validator(); err != nil {
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
	cfg.ZKTServer.URL.UserID = "http://money.fylos.cn:1234/selfservice/selfreport/"
	cfg.ZKTServer.URL.Login = "http://money.fylos.cn:1234/selfservice/login/" // XXX:host route-path split
	cfg.ZKTServer.URL.TimeTag = "http://money.fylos.cn:1234/grid/att/CardTimes/"
	cfg.XServer.Addr = "0.0.0.0:4567"
	cfg.XServer.File.Pid = filepath.Join(WorkDir, AppName+".pid")
	cfg.XServer.File.Log = filepath.Join(WorkDir, AppName+".log")
	cfg.XServer.DB.Type = "sqlite3"
	cfg.XServer.DB.Path = filepath.Join(WorkDir, "data", "data.db")
	cfg.XServer.ShortURL.Server.AppKey = "5db6aba18e676d1b43de23f6@79e2122d548ba64431f097e6c516774d"
	cfg.XServer.ShortURL.Server.ApiAddr = "http://api.suolink.cn/api.php"
	cfg.XServer.ShortURL.PrefixURL = "http://fylos.cn:4567/api/v1"
	cfg.XClient.Server.Addr = "http://127.0.0.1:4567"
	cfg.DelayWorkTime.Item = []struct {
		Time  string `json:"time"`
		Delay uint32 `json:"delay"`
	}{
		{Time: "20:00:00", Delay: 15},
		{Time: "21:00:00", Delay: 45},
		{Time: "22:00:00", Delay: 105},
		{Time: "23:00:00", Delay: 200},
	}

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

	dbtype := cfg.XServer.DB.Type
	switch dbtype {
	case "sqlite3":
		if cfg.XServer.DB.Path == "" {
			return errors.New("sqlite path is not configureated, it'll not use default value")
		}
	case "mysql":
	default:
		return fmt.Errorf("database type(%s) not supported", dbtype)
	}
	return nil
}

func (cfg config) LogName() string {
	return cfg.XServer.File.Log
}
