package iotctl

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

type Iotctl struct {
	Debug   bool
	BaseOrg string
	Router  *mux.Router
}

func (app *Iotctl) Initialize(BaseOrg string) {
	app.Debug = true
	app.BaseOrg = BaseOrg

	app.Router = nil

	log.Debug("Setting up routes")

	n := negroni.Classic()
	n.UseHandler(app.Router)
	http.ListenAndServe("localhost:8080", n)
}
