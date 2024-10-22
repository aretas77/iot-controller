package iotctl

import (
	"encoding/json"
	"fmt"
	"net/http"

	typesMQTT "github.com/aretas77/iot-controller/types/mqtt"
	"github.com/aretas77/iot-controller/web/iotctl/database/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (app *Iotctl) PublishAck(network string, mac string, location string,
	sendInterval float32) error {
	payload := typesMQTT.MessageAck{
		Network:      network,
		MAC:          mac,
		Location:     location,
		SendInterval: int(sendInterval),
	}

	resp, err := json.Marshal(payload)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"mac":      mac,
			"network":  network,
			"location": location,
		}).Error("failed to generate payload ack for node")
		return err
	}

	publishTopic := fmt.Sprintf("control/%s/%s/ack", network, mac)
	logrus.Infof("publish ack on %s", publishTopic)
	logrus.Debugf("ack payload = %s", string(resp))
	token := app.Plain.Client.Publish(publishTopic, 0, false, resp)
	if token.Error() != nil {
		logrus.WithError(token.Error()).WithFields(logrus.Fields{
			"mac":      mac,
			"network":  network,
			"location": location,
		}).Error("failed to publish ack")
		return token.Error()
	}

	return nil
}

// PublishStatsHades will send the given statistics to the Hades service.
func (app *Iotctl) PublishStatsHades(network string, mac string, stats typesMQTT.MessageStats) error {
	resp, err := json.Marshal(stats)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"mac":     mac,
			"network": network,
		}).Error("failed to generate payload stats for hades")
		return err
	}

	publishTopic := fmt.Sprintf("hades/%s/%s/statistics", network, mac)
	logrus.Infof("publish stats on %s", publishTopic)
	logrus.Debugf("stats payload = %s", string(resp))
	token := app.Plain.Client.Publish(publishTopic, 0, false, resp)
	if token.Error() != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"mac":     mac,
			"network": network,
		}).Error("failed to publish payload stats for hades")
		return err
	}

	return nil
}

// PublishUnregister will send a message to the device that it has been
// removed from the network and should return to sending 'Grettings' instead
// of its normal operation - sending statistics, etc.
// TODO: should probably add a Middleware for unregistering a `Node` and the
// middleware would call a publish method and await for response from the device.
func (app *Iotctl) PublishUnregister(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	setupHeader(&w)

	vars := mux.Vars(r)
	node := models.Node{}
	if err := app.sql.GormDb.First(&node, vars["id"]).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	payload := typesMQTT.MessageUnregister{}

	payload.Network = ""
	payload.MAC = node.Mac

	resp, err := json.Marshal(payload)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"mac":     payload.MAC,
			"network": payload.Network,
		}).Error("failed to generate payload unregister for node")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	publishTopic := fmt.Sprintf("control/%s/%s/unregister", payload.Network, payload.MAC)
	logrus.Infof("publish unregister on %s", publishTopic)
	logrus.Debugf("unregister payload = %s", string(resp))
	token := app.Plain.Client.Publish(publishTopic, 0, false, resp)
	if token.Error() != nil {
		logrus.WithError(token.Error()).WithFields(logrus.Fields{
			"mac":     payload.MAC,
			"network": payload.Network,
		}).Error("failed to publish unregister")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	next(w, r)
}
