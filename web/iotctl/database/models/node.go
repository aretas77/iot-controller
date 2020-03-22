package models

import (
	"errors"
	"time"

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
	IpAddress6          string    `json:"ipv6"`
	LastSentAck         time.Time `json:"last_sent_ack"`
	LastReceivedMessage time.Time `json:"last_received"`
	Status              Status    `json:"status" sql:"type:ENUM('acknowledged', 'registered')" gorm:"default:'acknowledged'"`
	AddedUsername       string    `json:"username"`
	BatteryMah          float32   `json:"battery_left_mah"`
	BatteryPercentage   float32   `json:"battery_left_per"`

	// a `Has One` relationship. Node 1 <-> 1 NodeSettings.
	// Node `Has One` Settings.
	SettingsID uint         `json:"-"`
	Settings   NodeSettings `json:"settings,omitempty"`

	// Belongs to only one Network and its ID is kept in `NetworkRefer`.
	NetworkRefer uint     `json:"-"`
	Network      *Network `json:"network,omitempty"`
}

// NodeStatisticsEntry is used to track various statistics of Node devices.
type NodeStatisticsEntry struct {
	ID           uint      `gorm:"primary_key"`
	CPULoad      int       `json:"cpu_load"`
	Temperature  float32   `json:"temperature"`
	TempReadTime time.Time `json:"temp_read_time"`

	// Refers to UserId to whom it belongs to.
	UserRefer uint `json:"user_refer"`
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
	Mac string `json:"mac" gorm:"not null"`

	// a `Has One` relationship. UnregisteredNode 0..* <-> 1 Network.
	// UnregisteredNode `Has One` Network.
	NetworkRefer uint     `json:"network_refer" gorm:"not null"`
	Network      *Network `json:"network,omitempty"`

	NodeRefer uint  `json:"-"`
	Node      *Node `json:"node,omitempty"`
}

// NodeSettings describes the settings that are used by the Node device.
type NodeSettings struct {
	ID           uint `gorm:"primary_key"`
	NodeID       uint `json:"-"`
	ReadInterval int  `json:"read_interval"`
	SendInterval int  `json:"send_interval"`
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
