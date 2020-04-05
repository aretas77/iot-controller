package models

import (
	"errors"
	"time"

	typesMQTT "github.com/aretas77/iot-controller/types/mqtt"
	"github.com/jinzhu/gorm"
)

var (
	// ErrNodeNilPass ...
	ErrNodeNilPass = errors.New("(nil) passed instead of (Node)")
	// ErrNodeNotFound ...
	ErrNodeNotFound = errors.New("(Node) not found")
	// ErrNodeSettingsNilPass ...
	ErrNodeSettingsNilPass = errors.New("(nil) passed instead of (NodeSettings)")
	// ErrNodeSettingsNotFound ...
	ErrNodeSettingsNotFound = errors.New("NodeSettings) not found")
)

type Status string

const (
	Acknowledged = "acknowledged"
	Registered   = "registered"
)

// Node describes information about a specific Node device.
type Node struct {
	gorm.Model
	Name                string    `json:"name"`
	Mac                 string    `json:"mac" gorm:"unique;not null"`
	Location            string    `json:"location"`
	IpAddress4          string    `json:"ipv4"`
	LastSentAck         time.Time `json:"last_sent_ack"`
	LastReceivedMessage time.Time `json:"last_received"`
	Status              Status    `json:"status" sql:"type:ENUM('acknowledged', 'registered')" gorm:"default:'acknowledged'"`
	AddedUsername       string    `json:"username"`
	BatteryMah          float32   `json:"battery_left_mah"`
	BatteryMahTotal     float32   `json:"battery_total_mah"`
	BatteryPercentage   float32   `json:"battery_left_per"`

	// a `Has One` relationship. Node 1 <-> 1 NodeSettings.
	// Node `Has One` Settings.
	SettingsID uint         `json:"-"`
	Settings   NodeSettings `json:"settings,omitempty"`

	// Belongs to only one Network and its ID is kept in `NetworkRefer`.
	NetworkRefer uint     `json:"-"`
	Network      *Network `json:"network,omitempty"`
}

// NodeStatisticsEntry is used to track various statistics of `Node` devices
// at a given point of time.
//
// CPULoad		- is measured in percentages (0;100]%
// Pressure		- is measured in Pa.
// Temperature	- is measured in Celsius.
// BatteryMah	- is measured in mAh.
type NodeStatisticsEntry struct {
	ID                uint      `gorm:"primary_key"`
	CPULoad           int       `json:"cpu_load"`
	Pressure          float32   `json:"pressure"`
	Temperature       float32   `json:"temperature"`
	TempReadTime      time.Time `json:"temp_read_time"`
	BatteryMah        float32   `json:"battery_left_mah"`
	BatteryPercentage float32   `json:"battery_left_per"`

	// Not used
	Consumed *typesMQTT.ConsumedFrame `json:"consumed_battery,omitempty"`

	// Refers to `Nodes` MAC address to whom it belongs to.
	NodeRefer string `json:"node_refer"`
	Node      *Node  `json:"node,omitempty"`
}

// UnregisteredNode is used to register node - User supplies MAC address of
// the Node and thus Node is Registered.
//
// UnregisteredNode is used when a User adds a Node with AddNode request.
// However, this Node is still not connected anywhere - it needs to be mapped
// with a Node device which is added by MQTT broker.
//
// Once a Node is registered - UnregisteredNode should be removed.
type UnregisteredNode struct {
	gorm.Model
	Mac           string `json:"mac" gorm:"not null"`
	AddedUsername string `json:"username" gorm:"not null"`
	Location      string `json:"location"`

	// a `Has One` relationship. UnregisteredNode 0..* <-> 1 Network.
	// UnregisteredNode `Has One` Network.
	NetworkRefer uint     `json:"network_refer" gorm:"not null"`
	Network      *Network `json:"network,omitempty"`

	NodeRefer uint  `json:"-"`
	Node      *Node `json:"node,omitempty"`
}

// NodeSettings describes the settings that are used by the Node device.
type NodeSettings struct {
	ID           uint   `gorm:"primary_key"`
	NodeID       uint   `json:"-"`
	ReadInterval int    `json:"read_interval"`
	SendInterval int    `json:"send_interval"`
	DataFileName string `json:"-"`
	DataLineFrom int    `json:"-"`
	DataLineTo   int    `json:"-"`
	DataCount    int    `json:"-"`
}

type NodeService interface {

	// Init will be used to initialize all needed information and migration
	// data.
	Init() error

	// Create should create a Node device in the specified Database.
	Create(n *Node) (*Node, error)

	// Get should return a Node by given ID.
	Get(nodeID string) (*Node, error)

	// All should return all Nodes.
	All() ([]Node, error)

	// Update should update the given Node with given values.
	Update(n *Node) (*Node, error)
}
