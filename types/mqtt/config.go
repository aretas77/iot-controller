package mqtt

// Broker struct shall contain the information about an MQTT broker.
type Broker struct {
	Type         string
	Server       string
	Username     string
	Password     string
	ClientId     string
	ProtocolVer  uint
	DeviceMac    string
	CleanSession bool
}

type BrokerDeviceInfo struct {
	DeviceMac       string
	TotalBatteryMah float32
	BatteryLeftMah  float32
}
