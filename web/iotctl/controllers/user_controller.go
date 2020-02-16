package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/aretas77/iot-controller/web/iotctl/database"
	"github.com/gorilla/mux"
)

type UserController struct {
	TableName string
	Database  *db.Database
}

func (u *UserController) Init() error {

	return nil
}

func (u *UserController) migrateUserGorm() error {

	return nil
}

func (u *UserController) setupHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods",
		"POST, GET, OPTIONS, DELETE, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Access-Control-Allow-Origin")
}

func (u *UserController) GetUser(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {

	vars := mux.Vars(r)

	fmt.Println(vars)

	u.setupHeader(&w)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}
