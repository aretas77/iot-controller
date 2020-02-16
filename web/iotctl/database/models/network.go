package models

import (
	"errors"
)

var (
	// ErrNetworkNilPass ...
	ErrNetworkNilPass = errors.New("(nil) passed instead of (Network)")
	// ErrNetworkNotFound ...
	ErrNetworkNotFound = errors.New("(Network) not found")
)

type Network struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Users User   `json:"users" db:"users"`
	Nodes Node   `json:"nodes" db:"nodes"`
}

type NetworkService interface {

	// Init will be used to initialize all needed information and migration
	// data for Networks.
	Init() error

	// Create should create a new network
	Create(n *Network) (*Network, error)

	// AddUser should add a given user into a specified network.
	AddUser(user User, networkID string) error

	Get(nodeID string) (*Node, error)

	All() ([]Network, error)

	Update(n *Node) (*Node, error)
}
