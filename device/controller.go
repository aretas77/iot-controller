package device

import (
	"errors"
	"runtime"
	"sync"

	"github.com/aretas77/iot-controller/device/hal"
	"github.com/aretas77/iot-controller/types/mqtt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// DeviceController is the main struct for managing whole simulation.
type DeviceController struct {
	NoTLSBroker string
	ListHAL     []string

	// MQTT connections.
	Plain  mqtt.MQTTConnection
	Secure mqtt.MQTTConnection

	// MQTT topics.
	PlainTopics  []mqtt.TopicHandler
	SecureTopics []mqtt.TopicHandler

	// MQTT secret for authentication.
	mqttSecret string

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
			d.Plain.Client.Publish(packet.Topic, packet.QoS, false, packet.Payload)
		}
	}
}

func (d *DeviceController) Init(host string) error {
	d.mqttQueue = make(chan Message, 10)
	d.broadcast = make(map[string]chan Message)
	d.wg = sync.WaitGroup{}

	// Initialize HALs
	//d.ListHAL = append(d.ListHAL, "esp32")

	d.Plain.Options = &MQTT.ClientOptions{}
	d.Plain.Options.SetProtocolVersion(3)
	d.Plain.Options.SetClientID("device-simulator")
	d.Plain.Options.SetCleanSession(true)
	d.Plain.Options.SetOrderMatters(true)
	d.Plain.Options.SetAutoReconnect(true)
	d.Plain.Options.SetUsername("devices")
	d.Plain.Options.SetPassword("test")
	d.Plain.Options.AddBroker(host)
	d.Plain.Client = MQTT.NewClient(d.Plain.Options)
	token := d.Plain.Client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if err := d.subscribeDevicePlainTopics(); err != nil {
		return errors.New("Failed to subscribe plain device topics")
	}

	return nil
}

func (d *DeviceController) Start(stop chan bool, devs []DeviceInfo) error {
	logrus.Info("starting device simulation")

	// Channel for controlling when to stop the working nodes.
	exit := make(chan struct{})

	for _, dev := range devs {
		d.broadcast[dev.MAC] = make(chan Message, 10)
		tempDevice := &NodeDevice{
			Name:           dev.Name,
			Mac:            dev.MAC,
			Network:        dev.Network,
			Send:           d.mqttQueue,
			Receive:        d.broadcast[dev.MAC],
			Location:       "",
			IpAddress4:     "",
			Status:         NodeDeviceNew,
			StatisticsFile: dev.Statistics,

			// Lets give each NodeDevice a reference to the main WorkGroup and
			// an exit channel.
			Wg:   &d.wg,
			Stop: exit,

			// Values derived from the model and adjustable by S-MQTT.
			ReadInterval: 0,
			SendInterval: 0,

			BatteryMah: dev.BatteryMah,
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
