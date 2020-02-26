package iotctl

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func (app *Iotctl) httpSetup() error {

	// Setup controllers
	app.setupRoutes()

	n := negroni.Classic()
	n.UseHandler(app.Router)

	app.httpServer = &http.Server{
		Addr:    app.ListenAddress,
		Handler: n,
	}

	app.Controller.NodeCtl.Hello()
	go func() {
		app.wg.Add(1)
		if err := app.httpServer.ListenAndServe(); err != nil {
			logrus.Fatal("Error while starting http server")
		}
		logrus.Infof("stopping http interface")
		app.wg.Done()
	}()

	logrus.Infof("Started http interface: %s", app.ListenAddress)
	return nil
}

// setupRoutes will prepare all endpoints for http server.
func (app *Iotctl) setupRoutes() {
	logrus.Debug("Setting up Routes")
	app.Router = mux.NewRouter()

	app.setupNode()
	app.setupUser()
}

func (app *Iotctl) setupNode() {

	logrus.Debugf("Attaching negroni handler on GET %s", "/nodes")
	app.Router.Handle(
		"/nodes",
		negroni.New(
			negroni.HandlerFunc(app.Controller.NodeCtl.GetNodes),
		)).Methods("GET")

	logrus.Debugf("Attaching negroni handler on GET %s", "/nodes/{id}")
	app.Router.Handle(
		"/nodes/{id}",
		negroni.New(
			negroni.HandlerFunc(app.Controller.NodeCtl.GetNode),
		)).Methods("GET")

	logrus.Debugf("Attaching negroni handler on POST %s", "/nodes")
	app.Router.Handle(
		"/nodes",
		negroni.New(
			negroni.HandlerFunc(app.Controller.NodeCtl.AddNode),
		)).Methods("POST")
}

func (app *Iotctl) setupUser() {
	logrus.Debugf("Attaching negroni handler on GET %s", "/users")
	app.Router.Handle(
		"/users",
		negroni.New(
			negroni.HandlerFunc(app.Controller.UserCtl.GetUsers),
		)).Methods("GET")

	logrus.Debugf("Attaching negroni handler on GET %s", "/users/{id}")
	app.Router.Handle(
		"/users/{id}",
		negroni.New(
			negroni.HandlerFunc(app.Controller.UserCtl.GetUserById),
		)).Methods("GET")

	logrus.Debugf("Attaching negroni handler on POST %s", "/users")
	app.Router.Handle(
		"/users",
		negroni.New(
			nil,
		)).Methods("POST")

}
