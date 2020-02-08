package iotctl

import (
	"net/http"
	"os"
	"regexp"

	"github.com/aretas77/iot-controller/web/iotctl/controllers"
	db "github.com/aretas77/iot-controller/web/iotctl/database"
	"github.com/aretas77/iot-controller/web/iotctl/routers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

// Iotctl for main IoT controller settings and config. Handles
// various HTTP endpoints.
type Iotctl struct {
	router     *mux.Router
	options    *Options
	controller *controllers.ApiController
	database   *db.Database
}

// Options for the IoT controller.
type Options struct {
	ListenAddress string
	RoutePrefix   string
	CORSOrigin    *regexp.Regexp
	Debug         DebugInfo
}

// DebugInfo for debugging related information.
type DebugInfo struct {
	Level        log.Level
	ReportCaller bool // false by default
}

// Initialize should initialize all required struct's for
// iotctl.
func (app *Iotctl) Initialize(opts Options) {
	// Setup inner Options struct's
	opts.Debug.setupDebug()

	// Setup Iotctl struct
	app.options = &opts
	app.router = routers.Routes()

	// Setup database
	app.database = &db.Database{}
	app.database.Init()

	n := negroni.Classic()
	n.UseHandler(app.router)
	http.ListenAndServe(opts.ListenAddress, n)
}

// setupDebug should setup the debug information
func (dbg *DebugInfo) setupDebug() {
	dbg.setupLog()
}

func (dbg *DebugInfo) setupLog() {
	log.SetLevel(dbg.Level)

	// Output stdout instead of the default stderr.
	log.SetOutput(os.Stdout)

	// Add calling method as field.
	log.SetReportCaller(dbg.ReportCaller)
}
