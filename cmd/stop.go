package cmd

import (
	"encoding/json"
	"log"

	"github.com/leaftree/ctnotify/pkg/config"
	jsonresp "github.com/leaftree/k8shelper/pkg/resp"
	"github.com/leaftree/k8shelper/pkg/xhttp"
	"github.com/urfave/cli"
)

var Stop = cli.Command{
	Name:  "stop",
	Usage: "Stop server",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "restart, r",
			Usage: "restart server(todo)",
		},
	},
	Action: actionStop,
}

func actionStop(c *cli.Context) error {
	config.NewConfig()

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
