package hermes

import (
	typesMQTT "github.com/aretas77/iot-controller/types/mqtt"
	hermesmq "github.com/aretas77/paho.mqtt.golang"
)

// CreateHermesMQConfiguration constructs mqtt.ClientOptions from given a
// configuration data.
func CreateHermesMQConfiguration(options typesMQTT.Broker) (*hermesmq.ClientOptions, error) {
	opts := hermesmq.NewClientOptions()
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
