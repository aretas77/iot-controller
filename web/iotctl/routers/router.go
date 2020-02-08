package routers

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Routes() *mux.Router {
	log.Debug("Setting up routes")
	router := mux.NewRouter()

	// Setup routes
	router = SetNodeRoutes(router)

	return router
}
