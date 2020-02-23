package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

var (
	// ErrNetworkNilPass ...
	ErrNetworkNilPass = errors.New("(nil) passed instead of (Network)")
	// ErrNetworkNotFound ...
	ErrNetworkNotFound = errors.New("(Network) not found")
)

type Network struct {
	gorm.Model
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
