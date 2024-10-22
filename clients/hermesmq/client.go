package hermes

import (
	"log"
	"time"

	typesMQTT "github.com/aretas77/iot-controller/types/mqtt"
	hermesmq "github.com/aretas77/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// ClientCreator defines the function for creating an MQTT client.
type ClientCreator func(options typesMQTT.Broker) (hermesmq.Client, error)

// MQTTClient is an implementation of the interface MQTTClient which uses
// underlying Eclipse Paho library.
type MQTTClient struct {
	clientID string
	log      *logrus.Logger
	client   hermesmq.Client
}

func NewHermesMQClient(options typesMQTT.Broker) (typesMQTT.MQTTClient, error) {
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
	logrus.Info("connecting with HermesMQ")

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
		return hermesmq.ErrNotConnected
	}

	if tkn := c.client.Publish(topic, qos, false, payload); tkn.Error() != nil {
		return tkn.Error()
	}
	return nil
}

// Subscribe ...
func (c *MQTTClient) Subscribe(topic string, qos uint8, callback typesMQTT.CustomMessageHandler) error {
	if !c.IsConnected() {
		return hermesmq.ErrNotConnected
	}

	handler := func(cl hermesmq.Client, message hermesmq.Message) {
		if callback != nil {
			callback(typesMQTT.MessageDevice{
				Payload: message.Payload(),
				QoS:     message.Qos(),
				Topic:   message.Topic(),
			})
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
		return hermesmq.ErrNotConnected
	}
	c.client.Unsubscribe(topic)
	return nil
}

// SetErrorLog will set the output of the ERROR logger in MQTT.
func (c *MQTTClient) SetErrorLog(logger *log.Logger) {
	hermesmq.ERROR = logger
}

// SetWarnLog will set the output of the WARN logger in MQTT.
func (c *MQTTClient) SetWarnLog(logger *log.Logger) {
	hermesmq.WARN = logger
}

// SetDebugLog will set the output of the WARN logger in MQTT.
func (c *MQTTClient) SetDebugLog(logger *log.Logger) {
	hermesmq.DEBUG = logger
}

// DefaultClientCreator returns a default function for creating MQTT client.
func DefaultClientCreator() ClientCreator {
	return func(options typesMQTT.Broker) (hermesmq.Client, error) {
		clientOptions, err := CreateHermesMQConfiguration(options)
		if err != nil {
			return nil, err
		}

		return hermesmq.NewClient(clientOptions), nil
	}
}
