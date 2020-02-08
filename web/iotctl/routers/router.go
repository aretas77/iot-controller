package routers

import (
	"github.com/aretas77/iot-controller/web/iotctl/controllers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Routes(ctl *controllers.ApiController) *mux.Router {
	log.Debug("Setting up Routes")
	router := mux.NewRouter()

	// Setup routes
	router = SetNodeRoutes(router, ctl)

	return router
}
