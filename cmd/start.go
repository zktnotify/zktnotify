package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli"
	"github.com/zktnotify/zktnotify/models"
	"github.com/zktnotify/zktnotify/pkg/config"
	"github.com/zktnotify/zktnotify/pkg/service"
	"github.com/zktnotify/zktnotify/router"
)

var Start = cli.Command{
	Name:  "start",
	Usage: "Start supervisor and run in background",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "foreground, f",
			Usage: "start in foreground",
		},
		cli.StringFlag{
			Name:  "conf, c",
			Usage: "config file",
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
	ctx, canceled := context.WithCancel(context.Background())

	_, err := config.NewConfig(c.String("conf"))
	if err != nil {
		log.Println(err)
		exit(1)
	}

	logPath := config.Config.LogName()
	logFd, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("create file %s failed: %v", logPath, err)
	}

	if c.Bool("foreground") {
		config.NewConfig(c.String("conf"))
		models.NewEngine()
		service.Service(ctx)

		hdlr := router.NewApiMux()
		svr := &http.Server{
			Addr:         config.Config.XServer.Addr,
			WriteTimeout: time.Second * 4,
			Handler:      hdlr,
		}

		go func() {
			if err := svr.ListenAndServe(); err != nil {
				log.Println(err)
				canceled()
				exit(0)
			}
		}()
		log.Printf("server started, listening on %s\n", config.Config.XServer.Addr)
		catchExitSignal(syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
		canceled()
	} else {
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

func catchExitSignal(sigs ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sigs...)

	for sig := range ch {
		if sig == syscall.SIGHUP {
			continue
		}
		log.Printf("Got signal: %v, exit\n", sig)
		break
	}
}

func exit(sig int) {
	time.Sleep(time.Millisecond * 200)
	os.Exit(sig)
}
