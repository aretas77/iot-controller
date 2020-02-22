package controllers

import (
	"encoding/json"
	"net/http"

	db "github.com/aretas77/iot-controller/web/iotctl/database"
	models "github.com/aretas77/iot-controller/web/iotctl/database/models"
	mysql "github.com/aretas77/iot-controller/web/iotctl/database/mysql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type NodeController struct {
	TableName string
	Database  *db.Database

	// Nodes will be saved at MySQL database so just keep a pointer into
	// MySql struct for easier access.
	sql *mysql.MySql
}

func (n *NodeController) Init() error {
	if n.Database == nil {
		logrus.Error("NodeController: Database is nil!")
	}

	if n.Database.GetMySql() == nil {
		logrus.Error("NodeController: failed to get MySQL instance")
	} else {
		n.sql = n.Database.GetMySql()
	}

	n.migrateNodeGorm()

	return nil
}

func (n *NodeController) migrateNodeGorm() error {
	// Setup database for table creation
	n.sql.GormDb.DropTableIfExists(&models.Node{}, &models.NodeSettings{},
		&models.UnregisteredNode{})

	// Create tables
	n.sql.GormDb.CreateTable(&models.Node{}, &models.NodeSettings{},
		&models.UnregisteredNode{})

	// Add any required restrictions and foreign keys
	n.sql.GormDb.Model(&models.NodeSettings{}).AddForeignKey("node_refer", "nodes(id)", "RESTRICT", "RESTRICT")

	return nil
}

func (n *NodeController) setupHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods",
		"POST, GET, OPTIONS, DELETE, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Access-Control-Allow-Origin")
}

// GetNode should return a Node by its ID.
func (n *NodeController) GetNode(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {
	n.setupHeader(&w)

	vars := mux.Vars(r)
	node := models.Node{}
	n.sql.GormDb.First(&node, vars["id"])

	if node.Mac != "" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(node)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// GetNodes should return all registered Nodes.
func (n *NodeController) GetNodes(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {
	n.setupHeader(&w)

	nodes := []models.Node{}
	n.sql.GormDb.Find(&nodes)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nodes)
}

func (n *NodeController) AddNode(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {

	// Add node to unregistered nodes
	n.setupHeader(&w)
	decoder := json.NewDecoder(r.Body)

	logrus.Debugf("got request")
	var tmpNode models.Node
	if err := decoder.Decode(&tmpNode); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// Check if a Node with current MAC address exists.
		n.sql.GormDb.FirstOrCreate(&tmpNode, models.Node{Mac: tmpNode.Mac})
		w.WriteHeader(http.StatusCreated)
	}
}

// RegisterNode should add the Node to the specified network.
func (n *NodeController) RegisterNode(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {

	n.setupHeader(&w)
	decoder := json.NewDecoder(r.Body)

	var tmpNode models.UnregisteredNode
	if err := decoder.Decode(&tmpNode); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// Check if a given MAC address exists in Nodes.
	}
	// User inputs a MAC address of a device - pass UnregisteredNode object.

}

// UnRegisterNode should remove the Node from our network.
func (n *NodeController) UnregisterNode(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {

}
