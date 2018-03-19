package main

import (
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
	app.UsageText = "0-bundle [options] flist root"
	app.Version = "1.0"

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "env, e",
			Usage: "custom environemt variables",
		},
		cli.StringFlag{
			Name:  "redis, r",
			Value: "",
			Usage: "Redis server address for error reporting",
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

		return nil
	}
	app.Action = action

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
