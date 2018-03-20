package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("0-bundle")
)

func init() {
	logging.SetLevel(logging.INFO, "")
}

func main() {
	// 1- mount the flist
	// 2- we need to export some env variables
	// 3- how to run ? bash, etc.
	app := cli.NewApp()

	app.Name = "zbundle"
	app.Usage = "run mini environments from flist"
	app.UsageText = "zbundle [options] <flist>"
	app.Version = "1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "id, i",
			Usage: "[required] ID of the sandbox",
		},
		cli.StringSliceFlag{
			Name:  "env, e",
			Usage: "custom environemt variables",
		},
		cli.StringSliceFlag{
			Name:  "report, r",
			Usage: "report error back on failures to the given url",
		},
		cli.StringFlag{
			Name:  "storage, s",
			Value: "ardb://hub.gig.tech:16379",
			Usage: "storage url to use",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "run in debug mode",
		},
		cli.BoolFlag{
			Name:  "no-exit",
			Usage: "do not terminate (unmount the sandbox) after /etc/start exits",
		},
	}

	app.ArgsUsage = "flist"
	app.Before = func(ctx *cli.Context) error {
		formatter := logging.MustStringFormatter("%{time}: %{color}%{module} %{level:.1s} > %{message} %{color:reset}")
		logging.SetFormatter(formatter)
		if ctx.GlobalBool("debug") {
			logging.SetLevel(logging.DEBUG, "")
		} else {
			logging.SetLevel(logging.INFO, "")
		}

		//validate required inputs
		for _, key := range []string{"id"} {
			value := ctx.GlobalString(key)
			if value == "" {
				return fmt.Errorf("flag --%s is required", key)
			}
		}

		return nil
	}
	app.Action = action

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
