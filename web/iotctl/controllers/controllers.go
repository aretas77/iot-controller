package controllers

import (
	db "github.com/aretas77/iot-controller/web/iotctl/database"
	"github.com/sirupsen/logrus"
)

type ApiController struct {
	Version string
	NodeCtl *NodeController
}

func (api *ApiController) InitControllers(database *db.Database) error {

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

	return nil
}
