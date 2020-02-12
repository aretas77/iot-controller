package routers

import (
	"github.com/aretas77/iot-controller/web/iotctl/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetNodeRoutes should prepare all endpoints for `Node`
// manipulation
func SetNodeRoutes(router *mux.Router, ctl *controllers.ApiController) *mux.Router {

	router.Handle(
		"/nodes",
		negroni.New(
			negroni.HandlerFunc(ctl.NodeCtl.GetNode),
		)).Methods("GET")

	router.Handle(
		"/nodes/{id}",
		negroni.New(
			nil,
		)).Methods("GET")

	router.Handle(
		"/nodes",
		negroni.New(
			nil,
		)).Methods("POST")

	return router
}
