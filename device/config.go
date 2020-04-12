package device

import "github.com/aretas77/iot-controller/types/mqtt"

// DeviceInfo will define a device components from which it will
// be built.
type DeviceInfo struct {
	Name       string
	MAC        string
	Sensors    []string
	Network    string
	Battery    float32
	Interface  string
	Statistics string
	Ipaddress4 string
	Hermes     bool
}

type Config struct {
	Broker  mqtt.Broker
	Devices map[string]DeviceInfo `yaml:devices,omitempty`
}
