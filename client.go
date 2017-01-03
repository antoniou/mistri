package main

import (
	"fmt"
	"os"

	"github.com/antoniou/zero2Pipe/command"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Zero2Pipeline"
	app.Usage = ""
	app.Commands = []cli.Command{
		{
			Name:        "create",
			Usage:       "create a CD Pipeline",
			Description: "Create a CI/CD Pipeline",
			ArgsUsage:   "<pipeline>",
			Action: func(c *cli.Context) error {
				fmt.Println("Creating new pipeline: ", c.Args().First())
				cm := command.CreateCommand{}
				fmt.Println(cm.Run(c.Args()))
				return nil
			},
		},
	}

	app.Run(os.Args)
}
