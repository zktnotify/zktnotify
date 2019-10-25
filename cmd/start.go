package cmd

import (
	"context"
	"fylos/flog"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/leaftree/onoffice/models"
	"github.com/leaftree/onoffice/pkg/config"
	"github.com/leaftree/onoffice/pkg/service"
	"github.com/urfave/cli"
)

var Start = cli.Command{
	Name:  "start",
	Usage: "Start supervisor and run in background",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "foreground, f",
			Usage: "start in foreground",
		},
	},
	Action: actionStartServer,
}

func GoFunc(f func() error) chan error {
	ch := make(chan error)
	go func() {
		ch <- f()
	}()
	return ch
}

func actionStartServer(c *cli.Context) error {
	config.NewConfig()

	logPath := config.Config.LogName()
	logFd, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("create file %s failed: %v", logPath, err)
	}

	if c.Bool("foreground") {
		config.NewConfig()
		models.NewEngine()
		service.Service(context.Background())
	} else {
		flog.Info()
		cmd := exec.Command(os.Args[0], "start", "-f")
		cmd.Stdout = logFd
		cmd.Stderr = logFd
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		select {
		case err = <-GoFunc(cmd.Wait):
			log.Fatalf("server started failed, %v", err)
		case <-time.After(200 * time.Millisecond):
		}
	}
	return nil
}
