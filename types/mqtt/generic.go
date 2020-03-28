package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTConnection will represent a single MQTT connection with its options.
type MQTTConnection struct {
	Options *mqtt.ClientOptions
	Client  mqtt.Client
}

type TopicHandler struct {
	Topic   string
	Handler func(c mqtt.Client, msg mqtt.Message)
}

// MessageHandler acts as an overlay to the MessageHandler of any mqtt lib.
type MessageHandler func(c mqtt.Client, msg mqtt.Message)

// MQTTClient interface represents an underlying MQTT client implementation
// without regard to the used MQTT client library.
type MQTTClient interface {
	// IsConnected should return whether the client is connected to the broker.
	IsConnected() bool

	// Connect should connect to the MQTT broker.
	Connect() error

	// Disconnect should disconnect from the MQTT broker.
	Disconnect() error

	// Publish should publish the given payload to the given topic with a given QoS.
	Publish(topic string, qos uint8, payload interface{}) error

	// Subscribe should subscribe to the given topic with a given QoS and message handler.
	Subscribe(topic string, qos uint8, callback MessageHandler) error

	// Unsubscribe should unsubscribe from the given topic.
	Unsubscribe(topic string) error
}
