package iotctl

import (
	"net/http"
	"sync"

	"github.com/aretas77/iot-controller/web/iotctl/controllers"
	db "github.com/aretas77/iot-controller/web/iotctl/database"
	MQTT "github.com/eclipse/paho.mqtt.golang"
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

	die chan struct{}
	wg  sync.WaitGroup

	broker string

	// MQTT connections.
	Plain  MQTTConnection
	Secure MQTTConnection

	// MQTT topics.
	PlainTopics  []TopicHandler
	SecureTopics []TopicHandler

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

// MQTTConnection will represent a single MQTT connection with its options.
type MQTTConnection struct {
	Options *MQTT.ClientOptions
	Client  MQTT.Client
}

type TopicHandler struct {
	Topic   string
	Handler func(c MQTT.Client, msg MQTT.Message)
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

	return nil
}
