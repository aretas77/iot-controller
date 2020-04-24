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
	// BatteryConsumedSend indicates when a sensor was read + the values were
	// sent.
	BatteryConsumedSend = "send"
	// BatteryConsumedRead indicates when a sensor was read but there was no
	// send for the values.
	BatteryConsumedRead = "read"

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
// Statistics state:
//	In this state, after the model is received, the device periodically sends
//	sensor data until the device simulation is closed.
//
// NodeDevice is the main struct for device simulation.
type NodeDevice struct {
	devices.System

	LastSentGreeting    time.Time
	Time                time.Time
	DefaultReadInterval time.Duration

	// Statistics
	StatisticsFile  string
	SensorReadTimes int
	SendTimes       int

	wg  *sync.WaitGroup
	Hal hal.HAL // What Hardware Abstraction Layer is used

	// Channels for various communications between components.
	BatteryControl chan BatteryChangeInfo
	ReceivedAck    chan struct{}
	Unregister     chan struct{}
	Stop           chan struct{}
	Send           chan Message
	Receive        chan Message

	rwLock sync.RWMutex
}

// BatteryChangeInfo represents an event of how much battery was consumed
// and what action has consumed this amount.
type BatteryChangeInfo struct {
	consumed     float32
	consumedType string
}

func (n *NodeDevice) calculateBatteryPercentage() float32 {
	return (100 * n.System.CurrentBatteryMah) / n.System.BatteryMah
}

// Initialize will initialize the struct of NodeDevice.
func (n *NodeDevice) Initialize() error {
	// If HAL has failed to init - don't bother anymore, just return the error.
	if err := n.Hal.Initialize(); err != nil {
		return err
	}

	n.DefaultReadInterval = time.Duration(time.Second * 5)
	n.ReceivedAck = make(chan struct{})
	n.Unregister = make(chan struct{})
	n.BatteryControl = make(chan BatteryChangeInfo)
	return nil
}

// Start is essentially a 'state machine' for single NodeDevice which
// will transition between states.
// The NodeDevice will run as a goroutine which will have its own:
//	- ReceiveLoop
func (n *NodeDevice) Start() {

	// Initialize NodeDevice
	if err := n.Initialize(); err != nil {
		logrus.Error(err)
		return
	}
	defer n.Hal.PowerOff()

	ticker := time.NewTicker(n.DefaultReadInterval)
	// will handle the broadcasted messages from main device controller
	n.wg.Add(2)
	go n.ReceiveLoop()
	go n.MonitorDeviceLoop()

	// Send a greeting every N seconds defined by the ticker.
handshake:
	for {
		select {
		case <-n.ReceivedAck:
			logrus.Debugf("received ACK %s", n.System.Mac)
			break handshake
		case <-ticker.C:
			// publish a greeting
			n.PublishGreeting()
			n.LastSentGreeting = time.Now()
		case <-n.Stop:
			ticker.Stop()
			n.wg.Done()
			return
		}
	}

	ticker.Stop()

	// Won't be receiving anything on this channel.
	// close(n.ReceivedAck)

	// After a handshake is done, we can start tracking the time between sends
	// and how much energy was consumed during the given timeframe.
	//
	// NOTE: Each publish from this point should keep track of consumed
	// battery and how much time elapsed between various publish events.
	n.Time = time.Now()

	// Publish initial system information to verify that its correct and
	// don't wait for any response - continue to other state.
	n.PublishSystemData()

	statsTicker := time.NewTicker(n.DefaultReadInterval)
statistics:
	for {
		select {
		case <-statsTicker.C:
			logrus.Debugf("sending a statistic %s", n.System.Mac)
			n.PublishSensorData()

			// last sent
			n.Time = time.Now()
		case <-n.Unregister:
			goto handshake
		case <-n.Stop:
			n.wg.Done()
			break statistics
		}
	}

	logrus.Infof("Device stopped (%s)", n.System.Mac)
	return
}

// ReceiveLoop is individual for each of the device and will handle messages
// received on various topics.
//
// In short - this is internal device handlers.
func (n *NodeDevice) ReceiveLoop() {
	logrus.Debugf("receive loop running for %s", n.System.Mac)
	for {
		select {
		case msg := <-n.Receive:

			if msg.Topic == "ack" {
				// server acknowledged the device

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
				n.DefaultReadInterval = time.Minute * time.Duration(ack.SendInterval)
				logrus.Infof("Device (%s) status (%s) -> (%s), send stats every %v",
					n.System.Mac, NodeDeviceNew, n.System.Status, n.DefaultReadInterval)

				// proceed to the next state
				n.ReceivedAck <- struct{}{}
			} else if msg.Topic == "unregister" {
				// server wants to unregister the device

				unregister := mqtt.MessageUnregister{}
				if err := json.Unmarshal(msg.Payload, &unregister); err != nil {
					logrus.WithError(err).WithFields(logrus.Fields{
						"topic": msg.Topic,
						"msg":   msg.Payload,
					}).Error("failed to unmarshal unregister message")
				}

				logrus.Infof("Device (%s) status (%s) -> (%s)", n.System.Mac,
					n.System.Status, NodeDeviceNew)

				n.System.Status = NodeDeviceNew
				n.System.Network = "global"
				n.System.Location = ""

				n.Unregister <- struct{}{}
			} else if msg.Topic == "sent" {
				// notify Device monitor about a change in battery levels.
				n.SendTimes++
				n.BatteryControl <- BatteryChangeInfo{
					consumed:     0,
					consumedType: BatteryConsumedSend,
				}
			} else {
				logrus.Infof("%s <- %s. Payload:\n%s", n.System.Mac,
					msg.Topic, msg.Payload)
			}
		case <-n.Stop:
			return
		}
	}
}

// MonitorDeviceLoop will monitor the device battery levels and if it reaches
// 1 or less mAh - the device is stopped.
func (n *NodeDevice) MonitorDeviceLoop() {
	logrus.Debugf("(%s) monitor loop running", n.System.Mac)

	for {
		select {
		case control := <-n.BatteryControl:
			logrus.Debugf("(%s) received battery change event - %s",
				n.System.Mac, control.consumedType)

			// Action was made which was a success and so we can recalculate
			// battery levels.
			n.rwLock.Lock()

			switch control.consumedType {
			case BatteryConsumedRead:
				n.System.CurrentBatteryMah -= control.consumed
			case BatteryConsumedSend:
				n.System.CurrentBatteryMah -= (control.consumed +
					n.Hal.GetSendConsumed())
			}

			n.System.BatteryPercentage = n.calculateBatteryPercentage()

			logrus.Infof("(%s) battery change to %f", n.System.Mac,
				n.System.CurrentBatteryMah)
			if n.System.CurrentBatteryMah <= 1 {
				logrus.Infof("(%s) device battery level low - stop", n.System.Mac)
				n.Stop <- struct{}{}
			}

			n.rwLock.Unlock()
		case <-n.Stop:
			n.wg.Done()
		}
	}
}
