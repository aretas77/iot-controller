package iotctl

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// OnMessagePswRequest will handle password request by a node. This will
// public a psw into init/psw/response topic.
// Topic: init/psw/request
func (app *Iotctl) OnMessagePswRequest(client MQTT.Client, msg MQTT.Message) {
	logrus.Infof("plain message got on: %s", msg.Topic())
}
