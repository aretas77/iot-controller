package main

import (
	"os"
	"sort"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.NewApp()
	app.Name = "Device Simulator"
	app.Usage = "Simulating Devices"
	app.Description = "Simulates devices which are specified in a configuration file."
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "network, n",
			Usage:       "Specify a network where devices belong to",
			Value:       "",
			DefaultText: "global",
		},
		&cli.StringFlag{
			Name:  "server, s",
			Usage: "IP address of MQTT broker",
			Value: "localhost",
		},
		&cli.StringFlag{
			Name:  "data",
			Usage: "Sensor data for devices",
			Value: "",
		},
		&cli.BoolFlag{
			Name:  "debug, d",
			Usage: " Enable debugging",
		},
	}

	app.Commands = []*cli.Command{
		&cli.Command{
			Name:  "start",
			Usage: "Start the simulation",
			Flags: []cli.Flag{
				&cli.IntFlag{},
			},
			Action: func(c *cli.Context) error {
				// Read a config file

				// Read data file

				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	// Log Method Name
	//logrus.SetReportCaller(true)

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
