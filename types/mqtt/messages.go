package mqtt

// MessageGreeting is used by node devices when they first boot up. They send
// a so called Greeting/Hello to the server to indicate that they exist.
// TODO: currently, we need to manually set the Network name on the device,
// however, we need to make it so that network is returned by the server.
type MessageGreeting struct {
	MAC        string `json:"mac"`
	Name       string `json:"name"`
	IpAddress4 string `json:"ip_address4"`
}

// `MessageAck` is sent by the server after receiving a `MessageGreeting` from
// the device.
type MessageAck struct {
	MAC     string `json:"mac"`
	Network string `json:"network"`
}

type MessageStats struct {
	CPULoad     int     `json:"cpu_load"`
	BatteryLeft int     `json:"battery_left"`
	Temperature float32 `json:"temperature"`
	ReadTime    string  `json:"read_time"`
}
