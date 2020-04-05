package mqtt

import "time"

const (
	ConsumedStateHandshake  = "handshake"
	ConsumedStateStatistics = "statistics"
	ConsumedStateInit       = "init"
)

// MessageGreeting is used by node devices when they first boot up. They send
// a so called Greeting/Hello to the server to indicate that they exist.
// TODO: currently, we need to manually set the Network name on the device,
// however, we need to make it so that network is returned by the server.
type MessageGreeting struct {
	MAC        string    `json:"mac"`
	Name       string    `json:"name"`
	IpAddress4 string    `json:"ip_address4"`
	Sent       time.Time `json:"time_sent"`
}

// MessageAck is sent by the server after receiving a `MessageGreeting` from
// the device.
type MessageAck struct {
	MAC      string `json:"mac"`
	Network  string `json:"network"`
	Location string `json:"location"`
}

// MessageUnregister is sent by the server when a user wants to unregister the
// given `Node` from its network.
type MessageUnregister struct {
	MAC     string `json:"mac"`
	Network string `json:"network"`
}

// MessageEventSent is sent by the Hades daemone to the server to notify about
// a sent model to the device.
type MessageEventSent struct {
	Model    string `json:"model"`
	TimeSent string `json:"time_sent"`
}

// MessageStats is used for sending out various stats to the server.
type MessageStats struct {
	CPULoad           int           `json:"cpu_load"`
	BatteryLeft       float32       `json:"battery_left_mah"`
	BatteryPercentage float32       `json:"battery_left_per"`
	Temperature       float32       `json:"temperature"`
	Pressure          float32       `json:"pressure"`
	TempReadTime      time.Time     `json:"temp_read_time"`
	StatisticsCount   int           `json:"stats_line"`
	Consumed          ConsumedFrame `json:"consumed_battery"`
}

// ConsumedFrame is used to represent how much battery was consumed during
// the duration between the 'last send' and 'current send'.
type ConsumedFrame struct {
	Duration    time.Duration `json:"duration"`
	ConsumedMah float32       `json:"consumed_mah"`
	State       string        `json:"state"`
}
