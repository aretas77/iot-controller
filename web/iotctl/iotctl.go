package iotctl

import (
	"net/http"
	"sync"

	"github.com/aretas77/iot-controller/types/mqtt"
	"github.com/aretas77/iot-controller/web/iotctl/controllers"
	db "github.com/aretas77/iot-controller/web/iotctl/database"
	mysql "github.com/aretas77/iot-controller/web/iotctl/database/mysql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Iotctl for main IoT controller settings and config. Handles
// various HTTP endpoints.
type Iotctl struct {
	Router     *mux.Router
	Controller *controllers.ApiController
	Database   *db.Database
	UseGorm    bool
	sql        *mysql.MySql

	die chan struct{}
	wg  sync.WaitGroup

	// MQTT
	Plain        mqtt.MQTTConnection
	Secure       mqtt.MQTTConnection
	PlainTopics  []mqtt.TopicHandler
	SecureTopics []mqtt.TopicHandler
	mqttSecret   string
	broker       string

	// HTTP interface
	httpServer *http.Server

	ListenAddress string
	RoutePrefix   string

	Debug *DebugInfo
}

// DebugInfo for debugging related information.
type DebugInfo struct {
	Level        logrus.Level
	ReportCaller bool // false by default
}

func (app *Iotctl) Start() error {

	// setup http
	if err := app.httpSetup(); err != nil {
		logrus.Fatal("Failed to initialize HTTP server")
		return err
	}

	// Setup MQTT
	if err := app.ConnectMQTT(); err != nil {
		logrus.Fatal("Failed to initialize MQTT")
		return err
	}

	// TODO: Start a goroutine for monitoring Node's and UnregisteredNode's.

	return nil
}
