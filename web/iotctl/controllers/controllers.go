package controllers

import (
	"errors"

	db "github.com/aretas77/iot-controller/web/iotctl/database"
	"github.com/sirupsen/logrus"
)

var (
	ErrDatabaseNil = errors.New("ApiController Database is nil")
)

type ApiController struct {
	Version    string
	AuthCtl    *AuthController
	NodeCtl    *NodeController
	UserCtl    *UserController
	NetworkCtl *NetworkController
	EventCtl   *EventController
}

// Init should prepare and initialize all usable controllers with
// required options.
func (api *ApiController) Init(database *db.Database) error {
	logrus.Debug("Setting up Controllers")

	if database == nil {
		return ErrDatabaseNil
	}

	api.NodeCtl.Database = database
	api.UserCtl.Database = database
	api.NetworkCtl.Database = database
	api.AuthCtl.Database = database
	api.EventCtl.Database = database

	if err := api.UserCtl.Init(); err != nil {
		logrus.Error("Failed to initialize User Controller")
		return err
	}

	if err := api.NetworkCtl.Init(); err != nil {
		logrus.Error("Failed to initialize Network Controller")
		return err
	}

	if err := api.NodeCtl.Init(); err != nil {
		logrus.Error("Failed to initialize Node Controller")
		return err
	}

	if err := api.AuthCtl.Init(); err != nil {
		logrus.Error("Failed to initialize Auth Controller")
		return err
	}

	if err := api.EventCtl.Init(); err != nil {
		logrus.Error("Failed to initialize Event Controller")
		return err
	}

	return nil
}
