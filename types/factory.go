package types

import (
	"errors"
	"fmt"
	"strings"

	mqtt "github.com/aretas77/iot-controller/clients/mqtt"
	typesMQTT "github.com/aretas77/iot-controller/types/mqtt"
	"github.com/sirupsen/logrus"
)

const (
	// MQTT library implementation
	MQTT = "mqtt"

	// HermesMQ library implementation
	HermesMQ = "hermes"
)

// NewMqttClient will return
func NewMqttClient(config typesMQTT.Broker) (typesMQTT.MQTTClient, error) {
	logrus.Infof("attaching MQTT client type='%s'", config.Type)

	switch lowerType := strings.ToLower(config.Type); lowerType {
	case MQTT:
		return mqtt.NewMqttClient(config)
	case HermesMQ:
		return nil, nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown client type '%s' requested",
			config.Type))
	}

}
