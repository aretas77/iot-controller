package controllers

import (
	db "github.com/aretas77/iot-controller/web/iotctl/database"
	models "github.com/aretas77/iot-controller/web/iotctl/database/models"
	mysql "github.com/aretas77/iot-controller/web/iotctl/database/mysql"
	"github.com/sirupsen/logrus"
)

type NetworkController struct {
	TableName string
	Database  *db.Database

	// Networks will be saved at MySQL database so just keep a pointer into
	// MySql struct for easier access.
	sql *mysql.MySql
}

func (n *NetworkController) Init() error {
	if n.Database == nil {
		logrus.Error("NetworkController: Database is nil!")
	}

	if n.Database.GetMySql() == nil {
		logrus.Error("NetworkController: failed to get MySQL instance")
	} else {
		n.sql = n.Database.GetMySql()
	}

	n.migrateNetworkGorm()

	return nil
}

func (n *NetworkController) migrateNetworkGorm() error {
	globalNetwork := models.Network{
		Name:      "global",
		UserRefer: 1,
	}

	if n.sql.GormDb.NewRecord(globalNetwork) {
		n.sql.GormDb.Create(&globalNetwork)
	}

	return nil
}
