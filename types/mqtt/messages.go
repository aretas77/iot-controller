package mqtt

// MessageGreeting is used by node devices when they first boot up. They send
// a so called Greeting/Hello to the server to indicate that they exist.
// TODO: currently, we need to manually set the Network name on the device,
// however, we need to make it so that network is returned by the server.
type MessageGreeting struct {
	MAC     string `json:"mac"`
	Name    string `json:"name"`
	Network string `json:"network"`
}

// MessageAck is sent by the server after receiving a MessageGreetinge from
// the device.
type MessageAck struct {
	MAC     string `json:"mac"`
	Network string `json:"network"`
}

type MessageStats struct {
	CPULoad       int     `json:"cpu_load"`
	ConsumedPower int     `json:"consumed_power"`
	Temperature   float32 `json:"temperature"`
}
