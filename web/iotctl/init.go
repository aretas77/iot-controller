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

func (app *Iotctl) initServices() (err error) {

	if err = app.Controller.Init(app.Database); err != nil {
		logrus.Error("Failed to initialize Controllers")
		return
	}

	// XXX: this looks ugly, fuuuu
	if app.sql, err = app.Database.GetMySql(); err != nil {
		logrus.Error("Failed to assign Gorm Database to Iotctl struct")
		return
	}

	return
}
