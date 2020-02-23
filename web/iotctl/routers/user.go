package routers

import (
	"github.com/aretas77/iot-controller/web/iotctl/controllers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

// SetUserRoutes should prepare all endpoints for `Node`
// manipulation
func SetUserRoutes(router *mux.Router, ctl *controllers.ApiController) *mux.Router {

	logrus.Debugf("Attaching negroni handler on GET %s", "/users")
	router.Handle(
		"/users",
		negroni.New(
			negroni.HandlerFunc(ctl.UserCtl.GetUsers),
		)).Methods("GET")

	logrus.Debugf("Attaching negroni handler on GET %s", "/users/{id}")
	router.Handle(
		"/users/{id}",
		negroni.New(
			negroni.HandlerFunc(ctl.UserCtl.GetUserById),
		)).Methods("GET")

	logrus.Debugf("Attaching negroni handler on POST %s", "/users")
	router.Handle(
		"/users",
		negroni.New(
			nil,
		)).Methods("POST")

	return router
}
