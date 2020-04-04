package iotctl

import (
	"errors"

	"github.com/aretas77/iot-controller/types/mqtt"
	"github.com/sirupsen/logrus"
)

// subscribePlainTopics will subscribe to known topics using plain
// connection.
func (app *Iotctl) subscribePlainTopics() error {
	app.PlainTopics = []mqtt.TopicHandler{
		{"init/psw/request", app.OnMessagePswRequest},
		{"control/+/+/greeting", app.OnMessageGreeting},
		{"node/+/+/stats", app.OnMessageStats},
		{"node/+/+/system", app.OnMessageSystem},
		{"node/+/+/hades/event/+", app.OnMessageEvent},
	}

	if !app.Plain.Client.IsConnected() {
		return errors.New("failed to initialize plain connection")
	}

	for _, t := range app.PlainTopics {
		token := app.Plain.Client.Subscribe(t.Topic, 0, t.Handler)
		if token.WaitTimeout(60); token.Error() != nil {
			return token.Error()
		}
		logrus.Debugf("plain subscribed topic %s", t.Topic)
	}

	return nil
}

func (app *Iotctl) subscribeSecureTopics() error {
	app.SecureTopics = []mqtt.TopicHandler{
		{"node/+/+/model/update", nil},
		{"system/node/add", nil},
		{"system/node/delete", nil},
		{"system/node/update", nil},
		{"system/node/move", nil},
	}
	return nil
}

// unsubscribePlainTopics will unsubscribe topics that are used for plain
// connection.
func (app *Iotctl) unsubscribePlainTopics() {
	if !app.Plain.Client.IsConnected() {
		return
	}

	for _, t := range app.PlainTopics {
		app.Plain.Client.Unsubscribe(t.Topic)
		logrus.Debugf("plain unsubscribing: %s", t.Topic)
	}

}
