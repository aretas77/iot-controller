package device

import (
	"encoding/json"

	typesMQTT "github.com/aretas77/iot-controller/types/mqtt"
	"github.com/sirupsen/logrus"
)

// HandleAck will handle server Acknowledge message. It will parse the message
// and get the MAC address from it so it would know to which device to send
// the Ack.
//
// After pre-processing is done - pass the message to the device indicated
// by a MAC address.
func (d *DeviceController) HandleAck(msg typesMQTT.MessageDevice) {
	logrus.Infof("plain got message on: %s", msg.Topic)
	payload := typesMQTT.MessageAck{}

	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
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
		Payload: msg.Payload,
	}
}

// MessageHandler will handle the MQTT messages. When a message is received,
// the handler will broadcast the message to all Nodes.
// Nodes should only process their own messages.
func (d *DeviceController) HandleBroadcast(msg typesMQTT.MessageDevice) {
	for _, v := range d.broadcast {
		v <- Message{
			Topic:   msg.Topic,
			Payload: msg.Payload,
		}
	}
}

// HandleUnregister will handle the unregister messages sent by the server.
// When such message is received, the device should stop its current state
// and return to the `Greeting` state.
func (d *DeviceController) HandleUnregister(msg typesMQTT.MessageDevice) {
	logrus.Infof("plain got message on: %s", msg.Topic)
	payload := typesMQTT.MessageUnregister{}

	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"topic": msg.Topic,
			"msg":   msg.Payload,
		}).Error("failed to unmarshal ack message")
		return
	}

	// We know to which device the Unregister was sent to - send the Unregister
	// to the device via channel.
	d.broadcast[payload.MAC] <- Message{
		Topic:   "unregister",
		Payload: msg.Payload,
	}

}
