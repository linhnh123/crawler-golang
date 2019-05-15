package main

import (
	"crawler/consumer"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"s"},
			Action: func(c *cli.Context) error {
				return consumer.Crawler.Start()
			},
		},
	}

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "s")
	}

	app.Run(os.Args)
}
