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

func (n *NodeController) GetNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fmt.Println(vars)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}
