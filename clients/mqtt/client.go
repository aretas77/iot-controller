package mqtt

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"

	typesMQTT "github.com/aretas77/iot-controller/types/mqtt"
)

// ClientCreator defines the function for creating an MQTT client.
type ClientCreator func(options typesMQTT.Broker) (mqtt.Client, error)

// MQTTClient is an implementation of the interface MQTTClient which uses
// underlying Eclipse Paho library.
type MQTTClient struct {
	clientID string
	log      *logrus.Logger
	client   mqtt.Client
}

// NewMqttClient will create a new paho.mqtt client with given broker options.
func NewMqttClient(options typesMQTT.Broker) (typesMQTT.MQTTClient, error) {
	mqttClient, err := DefaultClientCreator()(options)
	if err != nil {
		return &MQTTClient{}, err
	}

	return &MQTTClient{
		client:   mqttClient,
		log:      nil,
		clientID: options.ClientId,
	}, nil
}

// IsConnected ...
func (c *MQTTClient) IsConnected() bool {
	if c.client == nil {
		return false
	}
	return c.client.IsConnected()
}

// Connect ...
func (c *MQTTClient) Connect() error {
	logrus.Info("connecting with pahoMQTT")

	token := c.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Disconnect will disconnect from the given MQTT server.
func (c *MQTTClient) Disconnect() error {
	if c.IsConnected() {
		c.client.Disconnect(1000)
		c.client = nil
	}
	return nil
}

// Publish ...
func (c *MQTTClient) Publish(topic string, qos uint8, payload interface{}) error {
	if !c.IsConnected() {
		return nil
		//return mqtt.ErrNotConnected
	}
	c.client.Publish(topic, qos, false, payload)
	return nil
}

// Subscribe ...
func (c *MQTTClient) Subscribe(topic string, qos uint8, callback typesMQTT.MessageHandler) error {
	if !c.IsConnected() {
		return mqtt.ErrNotConnected
	}

	handler := func(cl mqtt.Client, message mqtt.Message) {
		if callback != nil {
			callback(c.client, message)
		}
	}

	token := c.client.Subscribe(topic, qos, handler)
	if token.WaitTimeout(3 * time.Second); token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Unsubscribe ...
func (c *MQTTClient) Unsubscribe(topic string) error {
	if !c.IsConnected() {
		return mqtt.ErrNotConnected
	}
	c.client.Unsubscribe(topic)
	return nil
}

// DefaultClientCreator returns a default function for creating MQTT client.
func DefaultClientCreator() ClientCreator {
	return func(options typesMQTT.Broker) (mqtt.Client, error) {
		clientOptions, err := CreateMQTTClientConfiguration(options)
		if err != nil {
			return nil, err
		}

		return mqtt.NewClient(clientOptions), nil
	}
}
