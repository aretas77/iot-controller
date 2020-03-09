package device

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// HandleAck will handle server Acknowledge message. It will parse the message
// and get the MAC address from it so it would know to which device to send
// the Ack.
func (d *DeviceController) HandleAck(client MQTT.Client, msg MQTT.Message) {
	logrus.Infof("plain got message on: %s", msg.Topic())
}

// MessageHandler will handle the MQTT messages. When a message is received,
// the handler will broadcast the message to all Nodes.
// Nodes should only process their own messages.
func (d *DeviceController) MessageHandler(c MQTT.Client, msg MQTT.Message) {
	for _, v := range d.broadcast {
		v <- Message{
			Topic:   msg.Topic(),
			Payload: msg.Payload(),
		}
	}
}
