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

// Network will be a bridge struct between a User and Nodes in the current
// Network.
//  - `Network` can have zero or more nodes. tl;dr Has Many `Nodes`.
//  - `Network` can belong to only one user. tl;dr Belongs To `User`.
//  - `Network` can be accessed by many users. Not implemented yet.
type Network struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `json:"name"`
	UserRefer uint   `json:"user_refer"`
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
