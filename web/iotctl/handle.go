package iotctl

import (
	"encoding/json"

	"github.com/aretas77/iot-controller/types/mqtt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// OnMessagePswRequest will handle password request by a node. This will
// public a psw into init/psw/response topic.
// Topic: init/psw/request
func (app *Iotctl) OnMessagePswRequest(client MQTT.Client, msg MQTT.Message) {
	logrus.Infof("plain message got on: %s", msg.Topic())
}

func (app *Iotctl) OnMessageGreeting(client MQTT.Client, msg MQTT.Message) {
	logrus.Infof("plain got message on: %s", msg.Topic())
	payload := mqtt.MessageAck{}

	// check if a valid network ID is provided
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"topic": msg.Topic,
			"msg":   msg.Payload,
		}).Error("failed to unmarshal ack message")
		return
	}

	// if valid, add into the unregistered Node db and send ack
	app.PublishAck(payload.Network, payload.MAC)
}
