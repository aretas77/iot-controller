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
	NodeCtl    *NodeController
	UserCtl    *UserController
	NetworkCtl *NetworkController
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

	err := api.NodeCtl.Init()
	if err != nil {
		logrus.Error("Failed to initialize Node Controller")
		return err
	}

	err = api.UserCtl.Init()
	if err != nil {
		logrus.Error("Failed to initialize User Controller")
		return err
	}

	err = api.NetworkCtl.Init()
	if err != nil {
		logrus.Error("Failed to initialize Network Controller")
		return err
	}

	return nil
}
