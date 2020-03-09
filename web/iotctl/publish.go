package iotctl

import (
	"encoding/json"
	"fmt"

	"github.com/aretas77/iot-controller/types/mqtt"
	"github.com/sirupsen/logrus"
)

func (app *Iotctl) PublishAck(network string, mac string) error {
	payload := mqtt.MessageAck{}

	payload.Network = network
	payload.MAC = mac

	resp, err := json.Marshal(payload)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"mac":     mac,
			"network": network,
		}).Error("failed to generate payload ack for node")
		return err
	}

	publishTopic := fmt.Sprintf("control/%s/%s/ack", network, mac)
	logrus.Infof("publish ack on %s", publishTopic)
	logrus.Debugf("ack payload = %s", string(resp))
	token := app.Plain.Client.Publish(publishTopic, 0, false, resp)
	if token.Error() != nil {
		logrus.WithError(token.Error()).WithFields(logrus.Fields{
			"mac":     mac,
			"network": network,
		}).Error("failed to publish ack")
		return token.Error()
	}

	return nil
}
