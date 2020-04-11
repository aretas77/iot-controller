package device

import (
	typesMQTT "github.com/aretas77/iot-controller/types/mqtt"
	"github.com/sirupsen/logrus"
)

func (d *DeviceController) subscribeDevicePlainTopics() error {
	d.PlainTopics = []typesMQTT.TopicHandlerDevice{
		{"control/+/+/ack", d.HandleAck},
		{"control/+/+/unregister", d.HandleUnregister},
	}

	for _, t := range d.PlainTopics {
		d.PlainConnection.Subscribe(t.Topic, 0, t.Handler)
		logrus.Debugf("plain subscribed topic %s", t.Topic)
	}

	return nil
}

func (d *DeviceController) unsubscribeDevicePlainTopics() {
	if !d.PlainConnection.IsConnected() {
		return
	}

	for _, t := range d.PlainTopics {
		d.PlainConnection.Unsubscribe(t.Topic)
		logrus.Debugf("plain unsubscribing: %s", t.Topic)
	}
}
