package cmd

import (
	"encoding/json"
	"log"

	"github.com/urfave/cli"
	"github.com/zktnotify/zktnotify/pkg/config"
	jsonresp "github.com/zktnotify/zktnotify/pkg/resp"
	"github.com/zktnotify/zktnotify/pkg/xhttp"
)

var Stop = cli.Command{
	Name:  "stop",
	Usage: "Stop server",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "restart, r",
			Usage: "restart server(todo)",
		},
		cli.StringFlag{
			Name:  "conf, c",
			Usage: "config file",
		},
	},
	Action: actionStop,
}

func actionStop(c *cli.Context) error {
	config.NewConfig(c.String("conf"))

	url := config.Config.XClient.Server.Addr + "/api/v1/shutdown"
	resp, err := xhttp.Get(url)
	if err != nil {
		log.Println("stop server failed:", err)
		return err
	}

	r := jsonresp.JSONResponse{}
	if err := json.Unmarshal(resp, &r); err != nil {
		log.Println("stop server response error:", err)
		return err
	}

	if r.Status != 0 {
		if r.Message != "" {
			log.Println(r.Message)
		} else if r.Data != nil {
			log.Println(r.Data)
		}
		return nil
	}
	log.Println("server stopped")

	return nil
}
