package device

import (
	"errors"
	"sync"

	"github.com/aretas77/iot-controller/types/mqtt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// DeviceController is the main struct for managing whole simulation.
type DeviceController struct {
	Options     *MQTT.ClientOptions
	Client      MQTT.Client
	NoTLSBroker string

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
	kill      bool
	stop      chan bool
	wg        sync.WaitGroup
}

type Message struct {
	Node    string
	Topic   string
	QoS     byte
	Payload []byte
}

func (d *DeviceController) PublishLoop(stop chan bool) {
	for {
		select {
		case <-stop:
			logrus.Info("exiting PublishLoop")
			return
		case packet := <-d.mqttQueue:
			logrus.Infof("%s -> %s (len:%d)", packet.Node, packet.Topic, len(packet.Payload))
		}
	}
}

// MessageHandler will handle the MQTT messages. When a message is received,
// the handler will broadcast the message to all Nodes.
// Nodes should only process their own messages.
func (d *DeviceController) MessageHandler(c MQTT.Client, msg MQTT.Message) {
	for _, v := range d.broadcast {
		v <- Message{
			Topic:   msg.Topic(),
			Payload: msg.Payload(),
		}
	}
}

func (d *DeviceController) Init(host string) error {
	d.mqttQueue = make(chan Message, 10)
	d.broadcast = make(map[string]chan Message)
	d.wg = sync.WaitGroup{}

	d.Plain.Options = &MQTT.ClientOptions{}
	d.Plain.Options.SetProtocolVersion(3)
	d.Plain.Options.SetClientID("device-simulator")
	d.Plain.Options.SetCleanSession(true)
	d.Plain.Options.SetOrderMatters(true)
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
			Name:       dev.Name,
			Mac:        dev.MAC,
			Network:    dev.Network,
			Send:       d.mqttQueue,
			Receive:    d.broadcast[dev.MAC],
			Location:   "",
			IpAddress4: "",
			Status:     NodeDeviceNew,

			// Lets give each NodeDevice a reference to the main WorkGroup.
			Wg:   &d.wg,
			Stop: exit,

			// Values derived from the model and adjustable by S-MQTT.
			ReadInterval: 0,
			SendInterval: 0,
		}

		// One new device going ONLINE!
		d.wg.Add(1)

		tempDevice.Start()
	}

	go d.PublishLoop(stop)

	//<-stop

	close(exit)
	d.wg.Wait()

	logrus.Infof("device controller is ending")

	return nil
}
