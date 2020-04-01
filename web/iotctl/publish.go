package iotctl

import (
	"encoding/json"
	"fmt"

	typesMQTT "github.com/aretas77/iot-controller/types/mqtt"
	"github.com/sirupsen/logrus"
)

func (app *Iotctl) PublishAck(network string, mac string, location string) error {
	payload := typesMQTT.MessageAck{}

	payload.Network = network
	payload.MAC = mac
	payload.Location = location

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

	publishTopic := fmt.Sprintf("node/%s/%s/hades/statistics", network, mac)
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
