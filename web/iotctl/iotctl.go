package iotctl

import (
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

type Iotctl struct {
	router  *mux.Router
	options *Options
}

type Options struct {
	ListenAddress string
	RoutePrefix   string
	CORSOrigin    *regexp.Regexp
	Debug         DebugInfo
}

type DebugInfo struct {
	Level        log.Level
	ReportCaller bool // false by default
}

// Initialize should initialize all required struct's for
// iotctl.
func (app *Iotctl) Initialize(opts Options) {
	// Setup Iotctl struct
	app.router = nil
	app.options = &opts

	// Setup inner Options struct's
	opts.Debug.setupDebug()

	log.Debug("Setting up routes")
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
