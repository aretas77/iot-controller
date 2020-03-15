package iotctl

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func stripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}

func (app *Iotctl) userAuthHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
		logrus.Infof("Auth token = %s", token)
	}

	next(w, r)
}
