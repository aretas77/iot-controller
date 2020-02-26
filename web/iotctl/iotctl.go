package iotctl

import (
	"net/http"
	"os"
	"regexp"
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
	router     *mux.Router
	options    *Options
	controller *controllers.ApiController
	database   *db.Database
	broker     string
	useGorm    bool

	die chan struct{}
	wg  sync.WaitGroup

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
}

// Options for the IoT controller.
type Options struct {
	ListenAddress string
	RoutePrefix   string
	CORSOrigin    *regexp.Regexp
	Debug         DebugInfo
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

// Initialize should initialize all required struct's for
// iotctl.
func (app *Iotctl) Initialize(opts Options) {
	// Setup inner Options struct's
	opts.Debug.setupDebug()

	app = &Iotctl{
		options:    &opts,
		database:   &db.Database{},
		controller: &controllers.ApiController{},
		router:     nil,
		useGorm:    true,
	}

	// Setup database
	app.database.Init(app.useGorm)
	// Deal with cleaning up
	defer app.database.Close(app.useGorm)

	// Setup controllers
	app.controller.InitControllers(app.database)

	// Setup MQTT
	if err := app.ConnectMQTT(); err != nil {
		logrus.Fatal("Failed to initialize MQTT")
		return
	}

	// Initialize greeting queue with a given queue size.
	app.GreetingQueueInit(100)

	// Start a goroutine for handling the Greetings sent from a device.
	go app.greetingQueueLoop(app.die)

	// setup http
	if err := app.httpSetup(); err != nil {
		logrus.Fatal("Failed to initialize HTTP server")
		return
	}

}

// setupDebug should setup the debug information
func (dbg *DebugInfo) setupDebug() {
	dbg.setupLog()
}

func (dbg *DebugInfo) setupLog() {
	logrus.SetLevel(dbg.Level)

	// Output stdout instead of the default stderr.
	logrus.SetOutput(os.Stdout)

	// Add calling method as field.
	logrus.SetReportCaller(dbg.ReportCaller)
}
