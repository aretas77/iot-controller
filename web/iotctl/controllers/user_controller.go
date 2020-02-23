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

type UserController struct {
	TableName string
	Database  *db.Database

	// Users will be saved at MySQL database so just keep a pointer into
	// MySql struct for easier access.
	sql *mysql.MySql
}

func (u *UserController) Init() error {
	if u.Database == nil {
		logrus.Error("UserController: Database is nil!")
	}

	if u.Database.GetMySql() == nil {
		logrus.Error("UserController: failed to get MySQL instance")
	} else {
		u.sql = u.Database.GetMySql()
	}

	u.migrateUserGorm()

	return nil
}

func (u *UserController) migrateUserGorm() error {
	// Tables are created - create an admin.
	user := models.User{
		Username: "Superadmin",
		Password: "test",
		Email:    "superadmin@gmail.com",
		Role:     "admin",
	}

	if u.sql.GormDb.NewRecord(user) {
		u.sql.GormDb.Create(&user)
	}

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

func (u *UserController) GetUserById(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {
	u.setupHeader(&w)

	vars := mux.Vars(r)
	user := models.User{}
	u.sql.GormDb.First(&user, vars["id"])

	if user.Username != "" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// GetUsers should return all system users.
// Method:		GET
// Endpoint:	/users
func (u *UserController) GetUsers(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {
	u.setupHeader(&w)

	users := []models.User{}
	u.sql.GormDb.Find(&users)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
