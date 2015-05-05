package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "dockerexec"
	app.Version = "1"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "log-file", Value: "", Usage: "set the log file to output logs to"},
		cli.BoolFlag{Name: "debug", Usage: "enable debug output in the logs"},
	}

	app.Commands = []cli.Command{
		execCommand,
	}

	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		if path := context.GlobalString("log-file"); path != "" {
			f, err := os.Create(path)
			if err != nil {
				return err
			}
			log.SetOutput(f)
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
