package controllers

import (
	db "github.com/aretas77/iot-controller/web/iotctl/database"
	"github.com/sirupsen/logrus"
)

type ApiController struct {
	Version string
	NodeCtl *NodeController
	UserCtl *UserController
}

// InitControllers should prepare and initialize all usable controllers with
// required options.
func (api *ApiController) InitControllers(database *db.Database) error {
	logrus.Debug("Setting up Controllers")

	// Initialize Node controller
	api.NodeCtl = &NodeController{
		TableName: "node",
		Database:  database,
	}

	err := api.NodeCtl.Init()
	if err != nil {
		logrus.Error("Failed to initialize Node Controller")
		return err
	}

	// Initialize User controller
	api.UserCtl = &UserController{
		TableName: "users",
		Database:  database,
	}

	err = api.UserCtl.Init()
	if err != nil {
		logrus.Error("Failed to initialize User Controller")
		return err
	}

	return nil
}
