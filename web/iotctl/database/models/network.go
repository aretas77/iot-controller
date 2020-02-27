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

	// a `Has Many` relationship. Network 1 <-> 0..* Node.
	// Network `Has Many` Nodes.
	Nodes []Node `json:"nodes" gorm:"foreignkey:NetworkRefer"`
}

type NetworkService interface {

	// Create should create a new network which belongs to some User -
	// should be specified in `UserRefer`.
	Create(net *Network) (*Network, error)

	// AddNode should add a given node into a specified network.
	AddNode(node *Node, networkId string) error

	// Get should return the Network of given ID.
	Get(networkId string) (*Network, error)

	// All should return all Networks.
	All() ([]Network, error)

	// Update ...
	Update(net *Network) (*Network, error)
}
