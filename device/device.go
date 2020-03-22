package device

import (
	"encoding/json"
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
//	After Handshake is done, the device will wait for a model for its MQTT
//	library. The model should contain initial values for when to read and send
//	the statistics. When the model is received, go to the Publish mode.
//
// Publish state:
//	In this state, after the model is received, the device periodically sends
//	sensor data until the device simulation is closed.
//
// NodeDevice is the main struct for device simulation.
type NodeDevice struct {
	devices.System

	// How much energy was consumed in the given time frame.
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

	// will handle the broadcasted messages from main device controller
	go n.ReceiveLoop()

	// Send a greeting every N seconds defined by the ticker.
handshake:
	for {
		select {
		case <-n.ReceivedAck:
			logrus.Debugf("received ACK %s", n.System.Mac)
			break handshake
		case <-ticker.C:
			// publish a greeting
			logrus.Debugf("sending a greeting %s", n.System.Mac)
			n.PublishGreeting()
		case <-n.Stop:
			ticker.Stop()
			n.Wg.Done()
			return
		}
	}

	// After a handshake is done, we can start tracking the time between sends
	// and how much energy was consumed during the given timeframe.
	n.ConsumedTimeFrame.Duration = time.Since(time.Now())

	// Publish initial system information to verify that its correct.
	n.PublishSystemData()

init:
	for {
		select {
		case <-n.Stop:
			n.Wg.Done()
			break init
		}
	}

	logrus.Infof("Device stopped (%s)", n.System.Mac)
	return
}

// ReceiveLoop is individual for each of the device and will handle messages
// received on various topics.
func (n *NodeDevice) ReceiveLoop() {
	logrus.Debugf("receive loop running for %s", n.System.Mac)
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
				// as existing and is registered to the network.
				n.System.Status = NodeDeviceRegistered
				n.System.Network = ack.Network
				n.System.Location = ack.Location
				logrus.Infof("Device (%s) status (%s) -> (%s)", n.System.Mac,
					NodeDeviceNew, n.System.Status)

				n.ReceivedAck <- struct{}{}
			} else {
				logrus.Infof("%s <- %s. Payload:\n%s", n.System.Mac,
					msg.Topic, msg.Payload)
			}
		case <-n.Stop:
			return
		}
	}
}
