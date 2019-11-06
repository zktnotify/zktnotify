package main

import (
	"log"
	"os"

	"github.com/leaftree/ctnotify/cmd"
	"github.com/urfave/cli"
)

const (
	version = "v0.1-alpha"
	pkg     = "ctnoitfy"
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
}

func main() {

	app := cli.NewApp()
	app.Name = pkg
	app.Usage = "一个打卡消息推送服务，推送上班、下班打卡消息；并且过了下班时间后提示你记得打卡"
	app.Version = version
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "fylos",
			Email: "fyl.root@gmail.com",
		},
	}
	app.Before = func(ctx *cli.Context) error {
		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf, c",
			Usage: "config file",
		},
	}
	app.Commands = []cli.Command{
		cmd.Start,
		cmd.Stop,
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
