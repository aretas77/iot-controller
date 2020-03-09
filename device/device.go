package device

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/aretas77/iot-controller/device/hal"
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
	Name       string
	Mac        string
	Location   string
	IpAddress4 string
	Status     string
	Network    string

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

// Start is essentially a 'state machine' for single NodeDevice which
// will transition between states.
// The NodeDevice will run as a goroutine which will have its own:
//	- ReceiveLoop
func (n *NodeDevice) Start() {
	n.ReceivedAck = make(chan struct{})
	ticker := time.NewTicker(3 * time.Second)

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

	return
}

func (n *NodeDevice) ReceiveLoop() {
	logrus.Debugf("receive loop running for %s", n.Mac)
	for {
		select {
		case msg := <-n.Receive:
			logrus.Infof("%s <- %s. Node: %s. Payload:\n%s", n.Mac, msg.Topic, msg.Node,
				msg.Payload)

			n.ReceivedAck <- struct{}{}
		case <-n.Stop:
			return
		}
	}
}

// PublishGreeting prepares a greeting message which will be sent to
// IoT controller which will add it as an `acknowledged` Node.
func (n *NodeDevice) PublishGreeting() {
	payload, _ := json.Marshal(&mqtt.MessageGreeting{
		MAC:     n.Mac,
		Name:    n.Name,
		Network: n.Network,
	})

	// Send to the main MQTT send channel
	n.Send <- Message{
		Node:    "",
		Topic:   fmt.Sprintf("control/%s/%s/greeting", n.Network, n.Mac),
		QoS:     0,
		Payload: payload,
	}
	return
}

// PublishSensorData prepares sensor data which will be sent for
// Reinforcement Learning.
func (n *NodeDevice) PublishSensorData() {

}
