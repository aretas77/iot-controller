package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	//BaseURL = "localhost:8080"
	app := cli.NewApp()
	app.Name = "iot-controller"

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "Starts IoT Controller Backend",
			Action: func(c *cli.Context) error {
				return start(c)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func start(c *cli.Context) error {
	return nil
}
