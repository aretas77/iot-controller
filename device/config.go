package device

// DeviceInfo will define a device components from which it will
// be built.
type DeviceInfo struct {
	Name    string
	Sensors []string
}

type Config struct {
	Broker struct {
		Server string
	}

	Devices map[string]DeviceInfo `yaml:devices,omitempty`
}
