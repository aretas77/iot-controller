package controllers

import (
	"encoding/json"
	"fmt"
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
	n.sql.GormDb.DropTableIfExists(&models.Node{}, &models.NodeSettings{})
	n.sql.GormDb.Model(&models.NodeSettings{}).AddForeignKey("node_refer", "nodes(id)", "RESTRICT", "RESTRICT")
	n.sql.GormDb.CreateTable(&models.Node{}, &models.NodeSettings{})
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

	vars := mux.Vars(r)

	query := "SELECT * FROM nodes"

	n.Database.GetMySql().Query(query)

	fmt.Println(vars)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}

func (n *NodeController) GetNodes(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {

	query := "SELECT * FROM nodes"
	n.Database.GetMySql().Query(query)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}

func (n *NodeController) AddNode(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {

	// Add node to unregistered nodes
}

// RegisterNode should add the Node to the specified network.
func (n *NodeController) RegisterNode(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {

	// User inputs a MAC address of a device

}

// UnRegisterNode should remove the Node from our network.
func (n *NodeController) UnregisterNode(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {

}
