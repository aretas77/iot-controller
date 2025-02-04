package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	db "github.com/aretas77/iot-controller/web/iotctl/database"
	models "github.com/aretas77/iot-controller/web/iotctl/database/models"
	mysql "github.com/aretas77/iot-controller/web/iotctl/database/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type NetworkController struct {
	TableName string
	Database  *db.Database

	// Networks will be saved at MySQL database so just keep a pointer into
	// MySql struct for easier access.
	sql *mysql.MySql
}

func (n *NetworkController) Init() (err error) {
	if n.Database == nil {
		return errors.New("NetworkController: Database is nil!")
	}

	if n.sql, err = n.Database.GetMySql(); err != nil {
		logrus.Error("NetworkController: failed to get MySQL instance")
		return err
	}

	n.migrateNetworkGorm()
	logrus.Debug("Initialized NetworkController")
	return
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

func (n *NetworkController) RemoveNetwork(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {
	n.setupHeader(&w)

	vars := mux.Vars(r)
	err := n.sql.GormDb.Where("id = ?", vars["id"]).Delete(models.Network{}).Error
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetUnregisteredNodesByNetwork ...
// Endpoint: GET /networks/{network_name}/unregistered
func (n *NetworkController) GetUnregisteredNodesByNetwork(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {
	n.setupHeader(&w)

	vars := mux.Vars(r)
	network := models.Network{}
	err := n.sql.GormDb.Where("name = ?", vars["network_name"]).First(&network).Error
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	unregisteredNodes := []models.UnregisteredNode{}
	err = n.sql.GormDb.Where("network_refer = ?", network.ID).Find(&unregisteredNodes).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mapNodes := map[string][]models.UnregisteredNode{}
	mapNodes["unregistered"] = unregisteredNodes

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mapNodes)
}

// GetNetworkByUser will return Network list which belong to the user of
// specified ID.
// Endpoint: GET /users/{user_id}/networks
func (n *NetworkController) GetNetworkByUser(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {
	n.setupHeader(&w)

	vars := mux.Vars(r)
	networks := []models.Network{}

	// Query all networks who belong to the specified user. Return error code
	// if an error occured.
	err := n.sql.GormDb.Where("user_refer = ?", vars["user_id"]).Find(&networks).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(networks)
}

// GetNetwork will return a single Network specified by ID. The returned Network
// should consist of its Nodes and their NodeSettings.
// Endpoint: GET /networks/{id}
func (n *NetworkController) GetNetwork(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {
	n.setupHeader(&w)

	vars := mux.Vars(r)
	network := models.Network{}
	nodes := []models.Node{}
	unregisteredNodes := []models.UnregisteredNode{}

	// SELECT * FROM `networks`  WHERE (`networks`.`id` = '1') ORDER BY `networks`.`id` ASC LIMIT 1
	if err := n.sql.GormDb.First(&network, vars["id"]).Error; err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Load Settings of Node device.
	// SELECT * FROM `node_settings`  WHERE (`node_id` IN (1,2))
	err := n.sql.GormDb.Where("network_refer = ?", network.ID).Preload("Settings").Find(&nodes).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = n.sql.GormDb.Where("network_refer = ?", network.ID).Find(&unregisteredNodes).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Assign found nodes
	network.Nodes = nodes
	network.UnregisteredNodes = unregisteredNodes

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(network)
}
