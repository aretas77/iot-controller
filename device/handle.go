package device

import (
	"encoding/json"

	"github.com/aretas77/iot-controller/types/mqtt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// HandleAck will handle server Acknowledge message. It will parse the message
// and get the MAC address from it so it would know to which device to send
// the Ack.
//
// After pre-processing is done - pass the message to the device indicated
// by a MAC address.
func (d *DeviceController) HandleAck(client MQTT.Client, msg MQTT.Message) {
	logrus.Infof("plain got message on: %s", msg.Topic())
	payload := mqtt.MessageAck{}

	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"topic": msg.Topic,
			"msg":   msg.Payload,
		}).Error("failed to unmarshal ack message")
		return
	}

	// We know to which device the ACK was sent to - send the ACK to the device
	// via channel.
	d.broadcast[payload.MAC] <- Message{
		Topic:   "ack",
		Payload: msg.Payload(),
	}
}

// MessageHandler will handle the MQTT messages. When a message is received,
// the handler will broadcast the message to all Nodes.
// Nodes should only process their own messages.
func (d *DeviceController) HandleBroadcast(c MQTT.Client, msg MQTT.Message) {
	for _, v := range d.broadcast {
		v <- Message{
			Topic:   msg.Topic(),
			Payload: msg.Payload(),
		}
	}
}
