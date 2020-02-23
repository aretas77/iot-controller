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
	ErrNodeNotFound = errors.New("(Device) not found")
	// ErrNodeSettingsNilPass ...
	ErrNodeSettingsNilPass = errors.New("(nil) passed instead of (NodeSettings)")
	// ErrNodeSettingsNotFound ...
	ErrNodeSettingsNotFound = errors.New("NodeSettings) not found")
)

type Status string

const (
	Acknowledged = "acknowledged"
	Registered   = "registered"
	Unregistered = "unregistered"
)

// Node describes information about a specific Node device.
type Node struct {
	gorm.Model
	Name        string    `json:"name" db:"name"`
	Mac         string    `json:"mac" db:"mac"`
	Location    string    `json:"location" db:"location"`
	IpAddress4  string    `json:"ipv4" db:"ipv4"`
	IpAddress6  string    `json:"ipv6" db:"ipv6"`
	LastSentAck time.Time `json:"last_sent_ack" db:"last_sent_ack"`
	Status      Status    `json:"status" sql:"type:ENUM('acknowledged', 'registered')" gorm:"default:'acknowledged'"`

	// a `Has One` relationship. Node 1 <-> 1 NodeSettings
	SettingsID uint `json:"settings"`

	// a `Belongs To` relationship. Node 0..* <-> 1 Network
	Network      Network `gorm:"foreignkey:NetworkRefer"`
	NetworkRefer uint    `json:"network_refer"`
}

type UnregisteredNode struct {
	gorm.Model
	Mac string `json:"mac" db:"mac"`
}

// NodeSettings describes the settings that are used by the Node device.
type NodeSettings struct {
	gorm.Model
	ReadInterval int `json:"read_interval" db:"read_interval"`
	SendInterval int `json:"send_interval" db:"send_interval"`
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
