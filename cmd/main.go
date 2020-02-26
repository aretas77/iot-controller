package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aretas77/iot-controller/web/iotctl"
	"github.com/aretas77/iot-controller/web/iotctl/controllers"
	db "github.com/aretas77/iot-controller/web/iotctl/database"
	mysql "github.com/aretas77/iot-controller/web/iotctl/database/mysql"
	"github.com/sirupsen/logrus"
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
	BaseURL = "localhost:8081"
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

func start(c *cli.Context) (err error) {
	MySqlDb := &db.Database{
		UseGorm: true,
		MySql: &mysql.MySql{
			Username: "root",
			Password: "test",
			Server:   "root:test@tcp(172.18.0.2:3306)/iotctl?parseTime=true",
		},
	}

	if err = MySqlDb.Init(); err != nil {
		return err
	}

	app := iotctl.Iotctl{
		ListenAddress: BaseURL,
		Database:      MySqlDb,
		Controller: &controllers.ApiController{
			NodeCtl: &controllers.NodeController{
				TableName: "nodes",
				Database:  MySqlDb,
			},
			UserCtl: &controllers.UserController{
				TableName: "users",
				Database:  MySqlDb,
			},
			NetworkCtl: &controllers.NetworkController{
				TableName: "networks",
				Database:  MySqlDb,
			},
		},
		Debug: &iotctl.DebugInfo{
			Level:        logrus.DebugLevel,
			ReportCaller: false,
		},
	}

	// listen for SIGTERM or SIGSTOP signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGSTOP,
		syscall.SIGKILL)

	if err := app.Init(); err != nil {
		return err
	}

	if err := app.Start(); err != nil {
		return err
	}

	<-stop
	return nil
}
