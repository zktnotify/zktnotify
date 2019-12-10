package cmd

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/urfave/cli"
	"github.com/zktnotify/zktnotify/pkg/config"
	jsonresp "github.com/zktnotify/zktnotify/pkg/resp"
	"github.com/zktnotify/zktnotify/pkg/xhttp"
)

var Status = cli.Command{
	Name:  "status",
	Usage: "show status of server",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "conf, c",
			Usage: "config file",
		},
	},
	Action: actionStatus,
}

func actionStatus(c *cli.Context) error {
	config.NewConfig(false, c.String("conf"))

	serverHost := hostname()

	started, err := isServerStartup()
	if err != nil {
		log.Println("get server status failed:", err)
		return err
	}
	if started {
		log.Printf("server(%s) is started up\n", serverHost)
		return nil
	}
	log.Printf("server(%s) is not started now\n", serverHost)

	return nil
}

func isServerStartup() (bool, error) {
	url := config.Config.XClient.Server.Addr + "/api/v1/status"
	rep, err := xhttp.Get(url)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return false, nil
		}
		return false, err
	}

	r := jsonresp.JSONResponse{}
	if err := json.Unmarshal(rep, &r); err != nil {
		return false, err
	}

	if r.Status != 0 {
		if r.Message != "" {
			return false, errors.New(r.Message)
		}
		return false, errors.New("unkonw error")
	}

	return true, nil
}

func hostname() string {
	return config.Config.XClient.Server.Addr
}
