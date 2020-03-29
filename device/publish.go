package device

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/aretas77/iot-controller/types/devices"
	"github.com/aretas77/iot-controller/types/mqtt"
)

// PublishGreeting prepares a greeting message which will be sent to
// IoT controller which will add it as an `acknowledged` Node.
func (n *NodeDevice) PublishGreeting() {
	payload, _ := json.Marshal(&mqtt.MessageGreeting{
		MAC:        n.System.Mac,
		Name:       n.System.Name,
		IpAddress4: n.System.IpAddress4,
		Sent:       time.Now(),
	})

	// Send to the main MQTT send channel
	n.Send <- Message{
		Mac:     n.System.Mac,
		Topic:   fmt.Sprintf("control/%s/%s/greeting", n.System.Network, n.System.Mac),
		QoS:     0,
		Payload: payload,
	}
	return
}

// PublishSystemData prepares system information and sends to the server
// to verify it.
func (n *NodeDevice) PublishSystemData() {
	payload, _ := json.Marshal(&devices.System{
		Mac:               n.System.Mac,
		Name:              n.System.Name,
		Location:          n.System.Location,
		Network:           n.System.Network,
		IpAddress4:        n.System.IpAddress4,
		Status:            n.System.Status,
		BatteryMah:        n.System.BatteryMah,
		BatteryPercentage: n.System.BatteryPercentage,
	})

	n.Send <- Message{
		Mac:     n.System.Mac,
		Topic:   fmt.Sprintf("node/%s/%s/system", n.System.Network, n.System.Mac),
		QoS:     0,
		Payload: payload,
	}
}

// PublishSensorData prepares sensor data which will be sent for
// Reinforcement Learning.
func (n *NodeDevice) PublishSensorData() {
	//consumed, temperature := n.Hal.GetTemperature("bmp180")
	consumed, temperature, pressure := n.Hal.GetPressureTemperature("bmp180")
	n.System.CurrentBatteryMah -= consumed
	n.System.BatteryPercentage = n.calculateBatteryPercentage()

	//n.ConsumedTimeFrame.ConsumedMah += consumed
	//n.ConsumedTimeFrame.Duration = time.Since(n.Time)

	payload, _ := json.Marshal(&mqtt.MessageStats{
		CPULoad:      rand.Intn(100),
		BatteryLeft:  n.System.CurrentBatteryMah,
		Temperature:  temperature,
		Pressure:     pressure,
		TempReadTime: time.Now(),
		Consumed:     n.ConsumedTimeFrame,
	})

	// Send to the main MQTT send channel
	n.Send <- Message{
		Mac:     n.System.Mac,
		Topic:   fmt.Sprintf("node/%s/%s/stats", n.System.Network, n.System.Mac),
		QoS:     0,
		Payload: payload,
	}

	return
}
