package mqtt

import MQTT "github.com/eclipse/paho.mqtt.golang"

// MQTTConnection will represent a single MQTT connection with its options.
type MQTTConnection struct {
	Options *MQTT.ClientOptions
	Client  MQTT.Client
}

type TopicHandler struct {
	Topic   string
	Handler func(c MQTT.Client, msg MQTT.Message)
}
