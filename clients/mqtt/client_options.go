package mqtt

import (
	typesMQTT "github.com/aretas77/iot-controller/types/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// CreateMQTTClientConfiguration construct mqtt.ClientOptions from given
// configuration data.
func CreateMQTTClientConfiguration(options typesMQTT.Broker) (*mqtt.ClientOptions, error) {
	opts := mqtt.NewClientOptions()
	opts.SetProtocolVersion(options.ProtocolVer)
	opts.SetClientID(options.ClientId)
	opts.SetCleanSession(true)
	opts.SetOrderMatters(true)
	opts.SetAutoReconnect(true)
	opts.SetUsername(options.Username)
	opts.SetPassword(options.Password)
	opts.AddBroker(options.Server)

	return opts, nil
}
