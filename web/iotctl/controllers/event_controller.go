package controllers

import (
	"errors"
	"net/http"

	typesMQTT "github.com/aretas77/iot-controller/types/mqtt"
	db "github.com/aretas77/iot-controller/web/iotctl/database"
	mysql "github.com/aretas77/iot-controller/web/iotctl/database/mysql"
	"github.com/sirupsen/logrus"
)

type EventController struct {
	TableName string
	Database  *db.Database

	// Nodes will be saved at MySQL database so just keep a pointer into
	// MySql struct for easier access.
	sql   *mysql.MySql
	plain *typesMQTT.MQTTConnection
}

func (e *EventController) Init() (err error) {
	if e.Database == nil {
		return errors.New("EventController: Database is nil!")
	}

	if e.sql, err = e.Database.GetMySql(); err != nil {
		logrus.Error("EventController: failed to get MySQL instance")
		return err
	}

	e.migrateEventGorm()
	logrus.Debug("Initialized EventController")
	return
}

func (e *EventController) migrateEventGorm() error {
	return nil
}

func (n *EventController) setupHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods",
		"POST, GET, OPTIONS, DELETE, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Access-Control-Allow-Origin")
}
