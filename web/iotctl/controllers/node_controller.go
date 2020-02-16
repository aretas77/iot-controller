package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/aretas77/iot-controller/web/iotctl/database"
	"github.com/gorilla/mux"
)

type NodeController struct {
	TableName string
	Database  *db.Database
}

func (n *NodeController) Init() error {
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

// RegisterNode should add the Node to the specified network.
func (n *NodeController) RegisterNode(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {

}

// UnRegisterNode should remove the Node from our network.
func (n *NodeController) UnregisterNode(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {

}
