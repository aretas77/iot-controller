package iotctl

import (
	"net/http"

	"github.com/aretas77/iot-controller/utils"
	"github.com/aretas77/iot-controller/web/iotctl/controllers"
	"github.com/sirupsen/logrus"
)

func setupHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods",
		"POST, GET, OPTIONS, DELETE, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Access-Control-Allow-Origin")
}

// userAuthBearer will validate an incoming request with a JWT Authorization
// header. If the JWT is invalid - stop the request chain.
func (app *Iotctl) userAuthBearer(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	setupHeader(&w)

	if isOptionsRequest(r) {
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		logrus.Debug("No Authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := utils.StripBearerPrefixFromTokenString(authHeader)
	if err != nil {
		logrus.Debug("Failed to strip Bearer prefix")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if token != "" {
		//logrus.Debugf("Auth token = %s", token)
		status, err := controllers.CheckBearerToken(token)
		if status != http.StatusOK || err != nil {
			logrus.Debugf("Failed")
			w.WriteHeader(status)
			return
		}
	}

	logrus.Info("Authenticated")
	next(w, r)
}

func isOptionsRequest(r *http.Request) bool {
	return r.Method == http.MethodOptions
}
