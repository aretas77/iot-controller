package iotctl

import (
	"net/http"

	"github.com/aretas77/iot-controller/web/iotctl/routers"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func (app *Iotctl) httpSetup() error {

	app.router = routers.Routes(app.controller)

	n := negroni.Classic()
	n.UseHandler(app.router)

	app.httpServer = &http.Server{
		Addr:    app.options.ListenAddress,
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

	logrus.Infof("Started http interface: %s", app.options.ListenAddress)
	return nil
}
