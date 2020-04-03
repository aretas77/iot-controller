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
	logrus.Debug("Setting up routing")
	app.Router = mux.NewRouter()

	// Currently everything is as admin.

	app.setupNode()
	app.setupUser()
	app.setupNetwork()
}

// setupNetwork will setup `Network` related routes.
func (app *Iotctl) setupNetwork() {
	app.Router.Handle(
		"/users/{user_id}/networks",
		negroni.New(
			negroni.HandlerFunc(app.Controller.UserCtl.Index),
		)).Methods("OPTIONS")

	app.Router.Handle(
		"/networks/{id}",
		negroni.New(
			negroni.HandlerFunc(app.Controller.UserCtl.Index),
		)).Methods("OPTIONS")

	app.Router.Handle(
		"/networks",
		negroni.New(
			negroni.HandlerFunc(app.Controller.NetworkCtl.CreateNetwork),
		)).Methods("POST")

	app.Router.Handle(
		"/networks/{id}",
		negroni.New(
			negroni.HandlerFunc(app.Controller.NetworkCtl.GetNetwork),
		)).Methods("GET")

	app.Router.Handle(
		"/users/{user_id}/networks",
		negroni.New(
			negroni.HandlerFunc(app.userAuthBearer),
			negroni.HandlerFunc(app.Controller.NetworkCtl.GetNetworkByUser),
		)).Methods("GET")
}

func (app *Iotctl) setupNode() {
	app.Router.Handle(
		"/nodes",
		negroni.New(
			negroni.HandlerFunc(app.Controller.UserCtl.Index),
		)).Methods("OPTIONS")
	app.Router.Handle(
		"/nodes/{id}",
		negroni.New(
			negroni.HandlerFunc(app.Controller.UserCtl.Index),
		)).Methods("OPTIONS")
	app.Router.Handle(
		"/nodes/{id}/statistics",
		negroni.New(
			negroni.HandlerFunc(app.Controller.UserCtl.Index),
		)).Methods("OPTIONS")

	app.Router.Handle(
		"/nodes",
		negroni.New(
			negroni.HandlerFunc(app.userAuthBearer),
			negroni.HandlerFunc(app.Controller.NodeCtl.GetNodes),
		)).Methods("GET")

	app.Router.Handle(
		"/nodes/{id}",
		negroni.New(
			negroni.HandlerFunc(app.userAuthBearer),
			negroni.HandlerFunc(app.Controller.NodeCtl.GetNode),
		)).Methods("GET")

	app.Router.Handle(
		"/nodes/{id}",
		negroni.New(
			negroni.HandlerFunc(app.userAuthBearer),
			negroni.HandlerFunc(app.PublishUnregister),
			negroni.HandlerFunc(app.Controller.NodeCtl.UnregisterNode),
		)).Methods("DELETE")

	app.Router.Handle(
		"/nodes/{id}/statistics",
		negroni.New(
			negroni.HandlerFunc(app.userAuthBearer),
			negroni.HandlerFunc(app.Controller.NodeCtl.GetEntries),
		)).Methods("GET")

	app.Router.Handle(
		"/nodes",
		negroni.New(
			negroni.HandlerFunc(app.Controller.NodeCtl.RegisterNode),
		)).Methods("POST")
}

func (app *Iotctl) setupUser() {
	// Routes for OPTIONS method
	app.Router.Handle(
		"/login",
		negroni.New(
			negroni.HandlerFunc(app.Controller.UserCtl.Index),
		)).Methods("OPTIONS")
	app.Router.Handle(
		"/users",
		negroni.New(
			negroni.HandlerFunc(app.Controller.UserCtl.Index),
		)).Methods("OPTIONS")
	app.Router.Handle(
		"/users/check",
		negroni.New(
			negroni.HandlerFunc(app.Controller.UserCtl.Index),
		)).Methods("OPTIONS")

	// Other routes
	app.Router.Handle(
		"/users/check",
		negroni.New(
			negroni.HandlerFunc(app.userAuthBearer),
			negroni.HandlerFunc(app.Controller.AuthCtl.CheckUsersToken),
		)).Methods("GET")

	app.Router.Handle(
		"/login",
		negroni.New(
			negroni.HandlerFunc(app.Controller.AuthCtl.Login),
		)).Methods("POST")

	app.Router.Handle(
		"/logout",
		negroni.New(
			negroni.HandlerFunc(app.userAuthBearer),
			negroni.HandlerFunc(app.Controller.AuthCtl.Logout),
		)).Methods("POST")

	app.Router.Handle(
		"/users",
		negroni.New(
			negroni.HandlerFunc(app.Controller.UserCtl.GetUsers),
		)).Methods("GET")

	app.Router.Handle(
		"/users/{id}",
		negroni.New(
			negroni.HandlerFunc(app.userAuthBearer),
			negroni.HandlerFunc(app.Controller.UserCtl.GetUserById),
		)).Methods("GET")
}
