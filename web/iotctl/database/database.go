package database

import (
	mysql "github.com/aretas77/iot-controller/web/iotctl/database/mysql"
	log "github.com/sirupsen/logrus"
)

// DatabaseService should be an abstract interface for various databases.
// e.g. MySQL or InnoDB.
type DatabaseService interface {

	// Init should initialize the database service.
	Connect() error

	// Query should send a given query to the database service.
	Query(query string) error

	// Close should close the database service.
	Close() error
}

// Database struct should keep the pointers to the realized interfaces.
type Database struct {
	MySql *mysql.MySql
}

func (d *Database) Init() {
	d.MySql = &mysql.MySql{}
	d.MySql.Connect()
}

func (d *Database) GetMySQL() *mysql.MySql {
	log.Debug("got ")
	return d.MySql
}
