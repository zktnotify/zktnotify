package cmd

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/zktnotify/zktnotify/pkg/version"
)

var Version = cli.Command{
	Name:   "version",
	Usage:  "show version",
	Flags:  []cli.Flag{},
	Action: actionVersion,
}

func actionVersion(c *cli.Context) error {
	fmt.Println(version.FullVersion())
	return nil
}
