package cmd

import (
	"log"

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
	log.Println("TODO")
	return nil
}
