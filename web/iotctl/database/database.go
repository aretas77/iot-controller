package database

import (
	mysql "github.com/aretas77/iot-controller/web/iotctl/database/mysql"
	"github.com/sirupsen/logrus"
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
func (d *Database) Init(useGorm bool) error {

	// TODO: pass this through config file when initializing resources.
	d.MySql = &mysql.MySql{
		Username: "root",
		Password: "test",
		Server:   "root:test@tcp(172.18.0.2:3306)/iotctl?parseTime=true",
	}

	if useGorm {
		d.MySql.ConnectGorm()
		d.MySql.InitializeMigrationGorm()
	} else {
		d.MySql.Connect()
	}

	return nil
}

// Close should close initialized databases.
func (d *Database) Close(useGorm bool) {
	logrus.Debugf("Closing Database connection. useGorm = %t", useGorm)
	if useGorm {
		d.MySql.CloseGorm()
	} else {
		d.MySql.Close()
	}
}

// GetMySql should return an initialized MySql object.
// TODO: make a single method which can return a DB by a given key.
func (d *Database) GetMySql() *mysql.MySql {
	return d.MySql
}
