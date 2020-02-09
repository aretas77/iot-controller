package iotctl

import (
	"runtime"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func (app *Iotctl) onConnectPlain(client MQTT.Client) {
	logrus.Info("Connected to plain MQTT")
	if err := app.subscribePlainTopics(); err != nil {
		logrus.Error("Failed to subscribe plain topics")
	}
}

// ConnectPlain should connect and subscribe topics using plain connection.
func (app *Iotctl) ConnectPlain() error {
	app.Plain.Options = &MQTT.ClientOptions{}
	app.Plain.Options.AddBroker("tcp://" + "172.18.0.3" + ":1883")
	app.Plain.Options.SetCleanSession(true)
	app.Plain.Options.SetClientID("iotctl-plain")
	app.Plain.Options.SetPingTimeout(60)
	app.Plain.Options.SetAutoReconnect(true)
	app.Plain.Options.SetOnConnectHandler(app.onConnectPlain)
	app.Plain.Options.SetConnectionLostHandler(nil)

	app.Plain.Client = MQTT.NewClient(app.Plain.Options)
	token := app.Plain.Client.Connect()
	if token.WaitTimeout(30); token.Error() != nil {
		return token.Error()
	}
	logrus.Infof("Plain connection initialized")

	return nil
}

// ConnectSecure should connect and subscribe topics using secure connection.
func (app *Iotctl) ConnectSecure() error {

	return nil
}

func (app *Iotctl) ConnectMQTT() error {
	logrus.Debug("Using MQTT secret " + "secret")
	if err := app.ConnectPlain(); err != nil {
		logrus.Debug("Plain connection failed: " + err.Error())
		return err
	}
	if err := app.ConnectSecure(); err != nil {
		logrus.Debug("Secure connection failed: " + err.Error())
		return err
	}
	logrus.Infof("MQTT (routines:%d)", runtime.NumGoroutine())

	return nil
}
