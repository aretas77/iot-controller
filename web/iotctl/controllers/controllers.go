package controllers

import (
	db "github.com/aretas77/iot-controller/web/iotctl/database"
	"github.com/sirupsen/logrus"
)

type ApiController struct {
	Version    string
	NodeCtl    *NodeController
	UserCtl    *UserController
	NetworkCtl *NetworkController
}

// InitControllers should prepare and initialize all usable controllers with
// required options.
func (api *ApiController) InitControllers(database *db.Database) error {
	logrus.Debug("Setting up Controllers")

	api = &ApiController{
		// Initialize Node controller
		NodeCtl: &NodeController{
			TableName: "nodes",
			Database:  database,
		},
		// Initialize User controller
		UserCtl: &UserController{
			TableName: "users",
			Database:  database,
		},
		// Initialize Network controller
		NetworkCtl: &NetworkController{
			TableName: "networks",
			Database:  database,
		},
	}

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
