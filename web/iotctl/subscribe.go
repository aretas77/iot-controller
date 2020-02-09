package iotctl

import (
	"errors"

	"github.com/sirupsen/logrus"
)

// subscribePlainTopics will subscribe to known topics using plain
// connection.
func (app *Iotctl) subscribePlainTopics() error {
	app.PlainTopics = []TopicHandler{
		{"init/psw/request", app.OnMessagePswRequest},
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
