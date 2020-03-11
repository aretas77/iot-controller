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

	sql *mysql.MySql

	die chan struct{}
	wg  sync.WaitGroup

	broker string

	// MQTT connections.
	Plain  mqtt.MQTTConnection
	Secure mqtt.MQTTConnection

	// MQTT topics.
	PlainTopics  []mqtt.TopicHandler
	SecureTopics []mqtt.TopicHandler

	// MQTT secret for authentication.
	mqttSecret string

	// HTTP interface
	httpServer *http.Server

	// When a Node sends a greeting to our controller, we don't reply
	// immediatly - we store it in a queue and send them with a delay.
	greetingQueue    *greetingEngine
	greetingQueueCap int

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
		return nil
	}

	// Setup MQTT
	if err := app.ConnectMQTT(); err != nil {
		logrus.Fatal("Failed to initialize MQTT")
		return nil
	}

	// Start a goroutine for handling the Greetings sent from a device.
	go app.greetingQueueLoop(app.die)

	// TODO: Start a goroutine for monitoring Node's and UnregisteredNode's.

	return nil
}
