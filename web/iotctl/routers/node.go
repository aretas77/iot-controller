package routers

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetNodeRoutes should prepare all endpoints for `Node`
// manipulation
func SetNodeRoutes(router *mux.Router) *mux.Router {

	router.Handle(
		"/nodes",
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
