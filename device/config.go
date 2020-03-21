package device

// DeviceInfo will define a device components from which it will
// be built.
type DeviceInfo struct {
	Name       string
	MAC        string
	Sensors    []string
	Network    string
	BatteryMah int `yaml:battery`
	Interface  string
	Statistics string
}

type Config struct {
	Broker struct {
		Server   string
		Username string
		Password string
	}

	Devices map[string]DeviceInfo `yaml:devices,omitempty`
}
