package main

import (
	"os"

	"github.com/aretas77/iot-controller/web/iotctl"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	// GitCommit is the last git commit passed with `ldflags`
	GitCommit string

	// Date is build date passed with `ldflags`
	Date string

	// Host should be passed with `ldflags`
	Host string

	// BaseURL is the url for the app
	BaseURL string
)

func main() {
	BaseURL = "localhost:8080"
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
	app := iotctl.Iotctl{}
	opts := iotctl.Options{}

	opts.ListenAddress = BaseURL
	opts.Debug = iotctl.DebugInfo{
		Level:        log.DebugLevel,
		ReportCaller: false,
	}

	app.Initialize(opts)

	return nil
}
