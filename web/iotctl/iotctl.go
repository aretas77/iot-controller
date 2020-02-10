package iotctl

import (
	"net/http"
	"os"
	"regexp"

	"github.com/aretas77/iot-controller/web/iotctl/controllers"
	db "github.com/aretas77/iot-controller/web/iotctl/database"
	"github.com/aretas77/iot-controller/web/iotctl/routers"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

// Iotctl for main IoT controller settings and config. Handles
// various HTTP endpoints.
type Iotctl struct {
	router     *mux.Router
	options    *Options
	controller *controllers.ApiController
	database   *db.Database
	broker     string

	// MQTT connections.
	Plain  MQTTConnection
	Secure MQTTConnection

	// MQTT topics.
	PlainTopics  []TopicHandler
	SecureTopics []TopicHandler

	// MQTT secret for authentication.
	mqttSecret string

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
	}

	// Setup database
	app.database.Init()

	// Setup controllers
	app.controller.InitControllers(app.database)

	// Setup routers
	app.router = routers.Routes(app.controller)

	// Setup MQTT
	if err := app.ConnectMQTT(); err != nil {
		logrus.Fatal("Failed to initialize MQTT")
		return
	}

	n := negroni.Classic()
	n.UseHandler(app.router)

	logrus.Debug("Listening to.. " + opts.ListenAddress)
	err := http.ListenAndServe(opts.ListenAddress, n)
	if err != nil {
		panic(err.Error())
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
