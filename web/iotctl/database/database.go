package database

import (
	mysql "github.com/aretas77/iot-controller/web/iotctl/database/mysql"
)

// DatabaseService should be an abstract interface for various databases.
// e.g. MySQL or InnoDB.
type DatabaseService interface {

	// Connect should initialize the database service using general MySQL
	// golang driver.
	Connect() error

	// ConnectGorm should initialize the database service using GORM.
	ConnectGorm() error

	// Query should send a given query to the database service.
	Query(query string, args ...interface{}) error

	// Close should close the database service.
	Close() error
}

// Database struct should keep the pointers to the realized interfaces.
type Database struct {
	MySql *mysql.MySql
	url   string
}

// Init should initialize all used databases.
func (d *Database) Init(useGorm bool) {
	d.MySql = &mysql.MySql{}
	d.url = "root:test@tcp(172.18.0.2:3306)/iotctl"

	if useGorm {
		d.MySql.ConnectGorm(d.url)
	} else {
		d.MySql.Connect(d.url)
	}
}

func (d *Database) GetMySql() *mysql.MySql {
	return d.MySql
}
