package database

import (
	"errors"

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
	MySql   *mysql.MySql
	UseGorm bool
}

// Init should initialize all used databases.
func (d *Database) Init() error {

	// TODO: we receive a MySql database struct via creation of DatabaseService
	// and thus we don't need to know any details of underlying database. So,
	// need to make a better abstraction without using individual values
	// unless configure a specific service using a configuration file.
	if d.UseGorm {
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
func (d *Database) GetMySql() (*mysql.MySql, error) {
	if d.MySql != nil {
		return d.MySql, nil
	}

	return nil, errors.New("MySql database not initialized")
}
