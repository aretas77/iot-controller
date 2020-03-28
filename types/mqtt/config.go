package mqtt

// Broker struct shall contain the information about an MQTT broker.
type Broker struct {
	Type        string
	Server      string
	Username    string
	Password    string
	ClientId    string
	ProtocolVer uint
}
