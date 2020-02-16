package mysql

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var (
	// ErrOpenFailed ...
	ErrOpenFailed = errors.New("Failed to open mysql database")
	// ErrPingFailed ...
	ErrPingFailed = errors.New("Database has failed to respond")
)

type MySql struct {
	server   string
	username string
	password string
	db       *sql.DB
	gormDb   *gorm.DB
}

func (m *MySql) ConnectGorm(url string) (err error) {
	logrus.Debug("Setting up MySQL database using GORM")

	m.gormDb, err = gorm.Open("mysql", url)
	if err != nil {
		logrus.Error(ErrOpenFailed)
		panic(err.Error())
	}

	err = m.db.Ping()
	if err != nil {
		logrus.Error(ErrPingFailed)
		panic(err.Error())
	}

	// Get the generic database object sql.DB to use its functions
	m.db = m.gormDb.DB()

	logrus.Debug("Connected to MySQL at 172.18.0.1:3306 with GORM")
	return
}

func (m *MySql) Connect(url string) (err error) {
	logrus.Debug("Setting up MySQL database")

	// Open a connection to MySQL database located at a specific IP.
	// This only returns a handle for a database. The database/sql package
	// manages connections in the background and doesn't open them until
	// we need it.
	m.db, err = sql.Open("mysql", url)
	if err != nil {
		logrus.Error("Failed to open mysql database")
		panic(err.Error())
	}

	err = m.db.Ping()
	if err != nil {
		panic(err.Error())
	}

	logrus.Debug("Connected to MySQL at 172.18.0.1:3306")
	return
}

func (m *MySql) Query(query string, args ...interface{}) error {
	rows, err := m.db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   int64
			name string
		)
		if err := rows.Scan(&id, &name); err != nil {
			logrus.Fatal(err)
		}
		log.Printf("id %d name is %s\n", id, name)
	}

	return nil
}

func (m *MySql) Close() error {
	err := m.db.Close()
	if err != nil {
		logrus.Fatal("Failed to close a database")
		logrus.Fatal(err)
	}
	return nil
}
