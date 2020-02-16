package models

import (
	"errors"
	"time"
)

var (
	// ErrNodeNilPass ...
	ErrNodeNilPass = errors.New("(nil) passed instead of (Node)")
	// ErrNodeNotFound ...
	ErrNodeNotFound = errors.New("(Device) not found")
)

type Node struct {
	ID           int       `json:"_key" db:"_key"`
	Location     string    `json:"location" db:"location"`
	SendInterval int       `json:"send_interval" db:"send_interval"`
	IpAddress4   string    `json:"ipv4" db:"ipv4"`
	IpAddress6   string    `json:"ipv6" db:"ipv6"`
	LastUpdated  time.Time `json:"last_update" db:"last_update"`
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
