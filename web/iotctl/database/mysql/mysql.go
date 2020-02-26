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
	// ErrOpenFailed is an error when Open method has failed.
	ErrOpenFailed = errors.New("Failed to open mysql database")
	// ErrPingFailed is an error when a Ping to the database has failed.
	ErrPingFailed = errors.New("Database has failed to respond")
)

// MySql struct should contain all information regarding a given connection
// to one MySql database. Also, it should implement a `DatabaseService`
// interface from database.go
type MySql struct {
	Server   string
	Username string
	Password string

	// We keep database connections open through all lifetime of our application.
	// Otherwise, with frequent Opens and Closes we could experience poor reuse,
	// sharing of connections, etc. More info:
	// - http://go-database-sql.org/accessing.html
	Db     *sql.DB
	GormDb *gorm.DB
}

// Connect should connect to a MySQL database using database/sql driver
// and gets the handle to the database.
func (m *MySql) Connect() (err error) {
	logrus.Debug("Setting up MySQL database")

	// Open a connection to MySQL database located at a specific IP.
	// This only returns a handle for a database. The database/sql package
	// manages connections in the background and doesn't open them until
	// we need it.
	m.Db, err = sql.Open("mysql", m.Server)
	if err != nil {
		logrus.Error(ErrOpenFailed)
		panic(err.Error())
	}

	err = m.Db.Ping()
	if err != nil {
		logrus.Error(ErrPingFailed)
		panic(err.Error())
	}

	logrus.Infof("Connected to MySQL at %s", m.Server)
	return
}

func (m *MySql) Query(query string, args ...interface{}) error {
	rows, err := m.Db.Query(query)
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

func (m *MySql) CloseGorm() error {
	err := m.GormDb.Close()
	if err != nil {
		logrus.Fatal("Failed to close a database")
		logrus.Fatal(err)
	}
	return nil
}

func (m *MySql) Close() error {
	err := m.Db.Close()
	if err != nil {
		logrus.Fatal("Failed to close a database")
		logrus.Fatal(err)
	}
	return nil
}
