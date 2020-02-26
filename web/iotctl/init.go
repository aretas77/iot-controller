package iotctl

import (
	"os"

	"github.com/sirupsen/logrus"
)

func (app *Iotctl) Init() error {

	if err := app.initLoggers(); err != nil {
		return err
	}

	if err := app.initServices(); err != nil {
		return err
	}

	if err := app.initQueues(); err != nil {
		return err
	}

	return nil
}

func (app *Iotctl) initLoggers() error {
	logrus.SetLevel(app.Debug.Level)

	// Output stdout instead of the default stderr.
	logrus.SetOutput(os.Stdout)

	// Add calling method as field.
	logrus.SetReportCaller(app.Debug.ReportCaller)

	return nil
}

func (app *Iotctl) initQueues() error {

	// Initialize greeting queue with a given queue size.
	app.GreetingQueueInit(100)

	return nil
}

func (app *Iotctl) initServices() error {

	if err := app.Database.Init(app.UseGorm); err != nil {
		logrus.Error("Failed to Initialize Database")
		return err
	}

	if err := app.Controller.Init(app.Database); err != nil {
		logrus.Error("Failed to Initialize Controllers")
		return err
	}

	return nil
}
