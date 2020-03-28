package device

import (
	"errors"
	"runtime"
	"sync"

	"github.com/aretas77/iot-controller/device/hal"
	"github.com/aretas77/iot-controller/types/devices"
	typesMQTT "github.com/aretas77/iot-controller/types/mqtt"
	"github.com/sirupsen/logrus"
)

// DeviceController is the main struct for managing whole simulation.
type DeviceController struct {
	ListHAL []string

	// MQTT connections.
	PlainConnection typesMQTT.MQTTClient

	// MQTT topics.
	PlainTopics  []typesMQTT.TopicHandler
	SecureTopics []typesMQTT.TopicHandler

	// Is used for sending out messages to all devices
	broadcast map[string]chan Message
	mqttQueue chan Message
	stop      chan bool
	wg        sync.WaitGroup
}

type Message struct {
	Mac     string
	Topic   string
	QoS     byte
	Payload []byte
}

// PublishLoop collects data from devices using `mqttQueue` channel and
// published it on the MQTT client.
func (d *DeviceController) PublishLoop(stop chan bool) {
	for {
		select {
		case <-stop:
			// if stop is nil - will block forever.
			logrus.Info("exiting PublishLoop")
			return
		case packet := <-d.mqttQueue:
			logrus.Infof("%s -> %s (len:%d)", packet.Mac, packet.Topic, len(packet.Payload))
			d.PlainConnection.Publish(packet.Topic, packet.QoS, packet.Payload)
		}
	}
}

func (d *DeviceController) Init(broker typesMQTT.Broker) error {
	d.mqttQueue = make(chan Message, 10)
	d.broadcast = make(map[string]chan Message)
	d.wg = sync.WaitGroup{}

	if err := d.PlainConnection.Connect(); err != nil {
		return errors.New("failed to connect plain connection")
	}

	if err := d.subscribeDevicePlainTopics(); err != nil {
		return errors.New("failed to subscribe plain device topics")
	}

	return nil
}

/*
func (d *DeviceController) initHermisMqtt(broker Broker) error {
	hermes.WARN = log.New(os.Stdout, "", 0)

	d.PlainHermes.Client = hermes.NewClient(d.PlainHermes.Options)
	token := d.PlainHermes.Client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
*/

func (d *DeviceController) Start(stop chan bool, devs []DeviceInfo) error {
	logrus.Info("starting device simulation")

	// Channel for controlling when to stop the working nodes.
	exit := make(chan struct{})

	// Construct NodeDevice with given and default parameters and run.
	for _, dev := range devs {
		d.broadcast[dev.MAC] = make(chan Message, 10)
		tempDevice := &NodeDevice{
			System: devices.System{
				Name:              dev.Name,
				Mac:               dev.MAC,
				Network:           dev.Network,
				Location:          "",
				IpAddress4:        "192.168.1.1",
				Status:            NodeDeviceNew,
				BatteryMah:        dev.Battery,
				CurrentBatteryMah: dev.Battery,
				BatteryPercentage: 100,
			},
			Send:           d.mqttQueue,
			Receive:        d.broadcast[dev.MAC],
			StatisticsFile: dev.Statistics,

			// Lets give each NodeDevice a reference to the main WorkGroup and
			// an exit channel.
			Wg:   &d.wg,
			Stop: exit,

			// Values derived from the model and adjustable by S-MQTT.
			// Should be parsed from the model.
			ReadInterval: 0,
			SendInterval: 0,
		}

		// Need to set the Hardware Abstraction Layer interface for the device.
		// TODO: make this better, somehow.
		switch dev.Interface {
		case "esp32":
			logrus.Infof("setting interface (%s) for device (%s)",
				dev.Interface, dev.Name)
			tempDevice.Hal = &hal.ESP32{}
			break
		default:
			logrus.Infof("interface not found (%s) for device (%s)",
				dev.Interface, dev.Name)
			break
		}

		// One new device going ONLINE!
		d.wg.Add(1)

		// This will start the device in the first mode: Handshake mode.
		go tempDevice.Start()
	}

	// There should be a single instance of PublishLoop.
	go d.PublishLoop(stop)

	logrus.Infof("active goroutines: %d", runtime.NumGoroutine())

	// nil stop will block forever
	<-stop

	close(exit)
	d.wg.Wait()

	logrus.Infof("device controller is ending")
	return nil
}
