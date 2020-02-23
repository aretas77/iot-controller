package routers

import (
	"github.com/aretas77/iot-controller/web/iotctl/controllers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

// SetNodeRoutes should prepare all endpoints for `Node`
// manipulation
func SetNodeRoutes(router *mux.Router, ctl *controllers.ApiController) *mux.Router {

	logrus.Debugf("Attaching negroni handler on GET %s", "/nodes")
	router.Handle(
		"/nodes",
		negroni.New(
			negroni.HandlerFunc(ctl.NodeCtl.GetNode),
		)).Methods("GET")

	logrus.Debugf("Attaching negroni handler on GET %s", "/nodes/{id}")
	router.Handle(
		"/nodes/{id}",
		negroni.New(
			nil,
		)).Methods("GET")

	logrus.Debugf("Attaching negroni handler on POST %s", "/nodes")
	router.Handle(
		"/nodes",
		negroni.New(
			negroni.HandlerFunc(ctl.NodeCtl.AddNode),
		)).Methods("POST")

	return router
}
