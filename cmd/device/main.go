package main

import (
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/aretas77/iot-controller/device"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func main() {
	var data_file string

	app := cli.NewApp()
	app.Name = "Device Simulator"
	app.Usage = "Simulating Devices"
	app.Description = "Simulates devices which are specified in a configuration file."
	app.EnableBashCompletion = true
	app.Compiled = time.Now()
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "Aretas Paulauskas",
			Email: "aretas.pau@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "network, n",
			Usage:       "Specify a network where devices belong to",
			DefaultText: "global",
		},
		&cli.StringFlag{
			Name:  "server, s",
			Usage: "IP address of MQTT broker",
			Value: "localhost",
		},
		&cli.StringFlag{
			Name:        "file, f",
			Usage:       "Sensor data file for devices",
			Value:       "configs/device.yaml",
			Destination: &data_file,
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
				return start(c, data_file)
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

func start(c *cli.Context, filename string) error {

	config := device.Config{}

	yamlConfig, err := ioutil.ReadFile(filename)
	if err != nil {
		logrus.Fatalf("Failed to open file: %s", err)
	}

	// Parsing configuration file
	if err := yaml.Unmarshal(yamlConfig, &config); err != nil {
		logrus.Fatal("Failed to parse yaml config")
	}

	// Need to map devices from map to list - we don't need a map.
	devices := []device.DeviceInfo{}
	for _, dev := range config.Devices {
		devices = append(devices, device.DeviceInfo{
			Name:    dev.Name,
			Sensors: dev.Sensors,
			Network: dev.Network,
		})
		logrus.Infof("Adding a device: %v", dev)
	}

	controller := &device.DeviceController{}
	if err := controller.Init("tcp://172.18.0.3:1883"); err != nil {
		return err
	}

	if err = controller.Start(nil, devices); err != nil {
		return err
	}

	return nil
}
