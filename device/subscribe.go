package device

import (
	"time"

	"github.com/aretas77/iot-controller/types/mqtt"
	"github.com/sirupsen/logrus"
)

func (d *DeviceController) subscribeDevicePlainTopics() error {
	d.PlainTopics = []mqtt.TopicHandler{
		{"control/+/+/ack", d.HandleAck},
		{"test", d.HandleBroadcast},
	}

	for _, t := range d.PlainTopics {
		token := d.Plain.Client.Subscribe(t.Topic, 0, t.Handler)
		if token.WaitTimeout(3 * time.Second); token.Error() != nil {
			return token.Error()
		}
		logrus.Debugf("plain subscribed topic %s", t.Topic)
	}

	return nil
}

func (d *DeviceController) unsubscribeDevicePlainTopics() {
	if !d.Plain.Client.IsConnected() {
		return
	}

	for _, t := range d.PlainTopics {
		d.Plain.Client.Unsubscribe(t.Topic)
		logrus.Debugf("plain unsubscribing: %s", t.Topic)
	}
}
