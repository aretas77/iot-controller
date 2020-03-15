package iotctl

import (
	"net/http"
	"strings"

	"github.com/aretas77/iot-controller/web/iotctl/controllers"
	"github.com/sirupsen/logrus"
)

func stripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}

// userAuthBearer will validate an incoming request with a JWT Authorization
// header. If the JWT is invalid - stop the request chain.
func (app *Iotctl) userAuthBearer(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := stripBearerPrefixFromTokenString(authHeader)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if token != "" {
		logrus.Debugf("Auth token = %s", token)
		status, err := controllers.CheckBearerToken(token)
		if status != http.StatusOK || err != nil {
			w.WriteHeader(status)
			return
		}
	}

	next(w, r)
}
