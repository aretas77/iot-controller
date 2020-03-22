package device

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/aretas77/iot-controller/device/hal"
	"github.com/aretas77/iot-controller/types/devices"
	"github.com/aretas77/iot-controller/types/mqtt"
	"github.com/sirupsen/logrus"
)

const (
	NodeDeviceAcknowledged = "acknowledged"
	NodeDeviceRegistered   = "registered"
	NodeDeviceNew          = "new"
)

// Device states
// 1. Handshake state - sending a greeting message until ack is received.
// 2. Initial model download state - downloading a model from the server for library.
// 3. Publish state - sending sensor data.
//
// Handshake state:
//	In this state, the device will continuously send a Greeting into
//	IoT controller service. When controller receives the Greeting - it will
//	then send a ReceivedAck message which will indicate that the device
//	has been added into the database with its status as `acknowledged`.
//	The ReceivedAck message will transfer the device into another state.
//
// Initial model download state: (skip for now)
//  After Handshake is done, the device will wait for a model for its MQTT
//	library. When the model is sent go to the Publish mode.
//
// Publish state:
//	In this state, after the model is received, the device periodically sends
//	sensor data until the device simulation is closed.
//
// NodeDevice is the main struct for device simulation.
type NodeDevice struct {
	devices.System

	// The battery size of current device.
	BatteryMah        float32
	BatteryPercentage float32
	ConsumedTimeFrame mqtt.ConsumedFrame

	// Statistics
	StatisticsFile string

	LastSentGreeting time.Time
	ReadInterval     time.Duration
	SendInterval     time.Duration

	Wg  *sync.WaitGroup
	Hal hal.HAL // What Hardware Abstraction Layer is used

	ReceivedAck chan struct{}
	Stop        chan struct{}

	Send    chan Message
	Receive chan Message
}

// Initialize will initialize the struct of NodeDevice.
func (n *NodeDevice) Initialize() error {
	// If HAL has failed to init - don't bother anymore, just return the error.
	if err := n.Hal.Initialize(n.StatisticsFile); err != nil {
		return err
	}

	n.ReceivedAck = make(chan struct{})

	return nil
}

// Start is essentially a 'state machine' for single NodeDevice which
// will transition between states.
// The NodeDevice will run as a goroutine which will have its own:
//	- ReceiveLoop
func (n *NodeDevice) Start() {
	ticker := time.NewTicker(3 * time.Second)

	// Initialize NodeDevice
	if err := n.Initialize(); err != nil {
		logrus.Error(err)
		return
	}
	defer n.Hal.PowerOff()

	n.Hal.GetTemperature("bmp180")

	// will handle the broadcasted messages from main device controller
	go n.ReceiveLoop()

handshake:
	for {
		select {
		case <-n.ReceivedAck:
			logrus.Infof("received ACK %s", n.Mac)
			break handshake
		case <-ticker.C:
			// publish a greeting
			logrus.Infof("sending a greeting %s", n.Mac)
			n.PublishGreeting()
		case <-n.Stop:
			ticker.Stop()
			n.Wg.Done()
			return
		}
	}

init:
	for {
		select {
		case <-n.Stop:
			n.Wg.Done()
			break init
		}
	}

	logrus.Infof("Device stopped (%s)", n.Mac)
	return
}

func (n *NodeDevice) ReceiveLoop() {
	logrus.Debugf("receive loop running for %s", n.Mac)
	for {
		select {
		case msg := <-n.Receive:

			if msg.Topic == "ack" {
				ack := mqtt.MessageAck{}
				if err := json.Unmarshal(msg.Payload, &ack); err != nil {
					logrus.WithError(err).WithFields(logrus.Fields{
						"topic": msg.Topic,
						"msg":   msg.Payload,
					}).Error("failed to unmarshal ack message")
				}

				// We now know that the device is acknowledged by the server
				// as existing.
				n.Status = NodeDeviceRegistered
				n.Network = ack.Network
				logrus.Infof("Device (%s) status (%s) -> (%s)", n.Mac,
					NodeDeviceNew, n.Status)

				n.ReceivedAck <- struct{}{}
			} else {
				logrus.Infof("%s <- %s. Payload:\n%s", n.Mac, msg.Topic, msg.Payload)
			}
		case <-n.Stop:
			return
		}
	}
}

// PublishGreeting prepares a greeting message which will be sent to
// IoT controller which will add it as an `acknowledged` Node.
func (n *NodeDevice) PublishGreeting() {
	payload, _ := json.Marshal(&mqtt.MessageGreeting{
		MAC:        n.Mac,
		Name:       n.Name,
		IpAddress4: "172.16.0.5",
		Sent:       time.Now(),
	})

	// Send to the main MQTT send channel
	n.Send <- Message{
		Mac:     n.Mac,
		Topic:   fmt.Sprintf("control/%s/%s/greeting", n.Network, n.Mac),
		QoS:     0,
		Payload: payload,
	}
	return
}

// PublishSensorData prepares sensor data which will be sent for
// Reinforcement Learning.
func (n *NodeDevice) PublishSensorData() {
	consumed, temperature := n.Hal.GetTemperature("bmp180")
	n.BatteryMah -= consumed

	payload, _ := json.Marshal(&mqtt.MessageStats{
		BatteryLeft:  n.BatteryMah,
		Temperature:  temperature,
		TempReadTime: time.Now(),
	})

	// Send to the main MQTT send channel
	n.Send <- Message{
		Mac:     n.Mac,
		Topic:   fmt.Sprintf("node/%s/%s/stats", n.Network, n.Mac),
		QoS:     0,
		Payload: payload,
	}
	return
}
