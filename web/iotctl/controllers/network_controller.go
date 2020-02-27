package controllers

import (
	"encoding/json"
	"net/http"

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

	logrus.Debug("Initialized NetworkController")
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

func (n *NetworkController) setupHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods",
		"POST, GET, OPTIONS, DELETE, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Access-Control-Allow-Origin")
}

// AddNetwork will create a Network which will belong to the specified User.
// User should be parsed from JWT token by UI and sent with the POST request.
// Endpoint: POST /networks
func (n *NetworkController) CreateNetwork(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {

	n.setupHeader(&w)
	decoder := json.NewDecoder(r.Body)

	var network models.Network
	if err := decoder.Decode(&network); err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		if network.UserRefer <= 0 {
			logrus.Info("CreateNetwork: network.UserRefer not defined")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		n.sql.GormDb.Create(&network)
		w.WriteHeader(http.StatusCreated)
	}
}
