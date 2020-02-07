package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type MySql struct {
	server   string
	username string
	password string
	db       *sql.Conn
}

func (m *MySql) Connect() error {
	log.Debug("Setting up MySQL database")

	db, err := sql.Open("mysql", "root:test@tcp(172.18.0.1:3306)/iotctl")
	if err != nil {
		log.Error("Failed to open mysql database")
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return nil
}

func (m *MySql) Query(query string) error {

	return nil
}

func (m *MySql) Close() error {

	return nil
}
