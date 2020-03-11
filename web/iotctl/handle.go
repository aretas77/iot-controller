package iotctl

import (
	"encoding/json"
	"time"

	"github.com/aretas77/iot-controller/types/mqtt"
	"github.com/aretas77/iot-controller/web/iotctl/database/models"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// OnMessagePswRequest will handle password request by a node. This will
// public a psw into init/psw/response topic.
// Topic: init/psw/request
func (app *Iotctl) OnMessagePswRequest(client MQTT.Client, msg MQTT.Message) {
	logrus.Infof("plain message got on: %s", msg.Topic())
}

func (app *Iotctl) OnMessageGreeting(client MQTT.Client, msg MQTT.Message) {
	logrus.Infof("plain got message on: %s", msg.Topic())
	payload := mqtt.MessageAck{}

	// check if a valid network ID is provided
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"topic": msg.Topic,
			"msg":   msg.Payload,
		}).Error("failed to unmarshal ack message")
		return
	}

	// check if a an `UnregisteredNode` is created with this MAC and depending
	// on whether it exists or not, branch out into two stages:
	//	1. `UnregisteredNode` doesn't exist - this means that the user hasn't
	//		requested any device with such MAC so we should ignore this devices
	//		greeting.
	//	2. `UnregisteredNode` with MAC exists - we know that the user has
	//		requested for this Node - we can send the ACK to the device and
	//		create a `Node` with fields: MAC and Network.

	tmpNode := models.UnregisteredNode{}
	err := app.sql.GormDb.Where("mac = ?", payload.MAC).Find(&tmpNode).Error
	if gorm.IsRecordNotFoundError(err) || err != nil {
		logrus.Error(err)

		// Error or `UnregisteredNode` not found - ignore this Greeting.
		logrus.Infof("greeting for %s but there is no entry for it - skip",
			payload.MAC)
		return
	}

	// At this point we know that some user is waiting for this device to
	// connect - delete the current `UnregisteredNode` entry and update an
	// existing `Node` or create a new `Node` with state 'Registered'.

	// Maybe a Node with such MAC already exists?
	// TODO: don't allow such inconsistencies.
	/*
		node := models.Node{}
		err = app.sql.GormDb.Where("mac = ?", tmpNode.Mac).Find(&node).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			logrus.Error(err)
			return
		}
	*/

	// First check if such network still exists as it could be that the user
	// has deleted it.
	network := models.Network{}
	err = app.sql.GormDb.Where("id = ?", tmpNode.NetworkRefer).Find(&network).Error
	if err != nil {
		logrus.Error(err)
		logrus.Infof("Network (ID = %d) doesn't exist anymore - skip",
			tmpNode.NetworkRefer)
		return
	}

	// `Node` doesn't exist and `Network` exists, lets create a `Node` in
	// a given `Network`.
	settings := models.NodeSettings{
		ReadInterval: 10,
		SendInterval: 10,
	}
	app.sql.GormDb.Create(&settings)

	node := models.Node{
		Name:         "AckNode",
		Mac:          payload.MAC,
		Location:     "",
		IpAddress4:   "",
		IpAddress6:   "NA",
		LastSentAck:  time.Now(),
		Status:       models.Registered,
		SettingsID:   settings.ID,
		NetworkRefer: tmpNode.NetworkRefer,
	}
	app.sql.GormDb.Create(&node)

	app.sql.GormDb.Model(&settings).Update("node_refer", node.ID)
	app.sql.GormDb.Delete(&tmpNode)

	logrus.Infof("Created a new Node (ID = %d, MAC = %s)", node.ID, node.Mac)

	// We can publish ack when we have done the following:
	//	1. An `UnregisteredNode` existed for this particular device node.
	//	2. We have created a new `Node` for this device with its status
	//	   as 'Registered'.
	app.PublishAck(network.Name, payload.MAC)
	return

	/*
		// `Node` exists.
		if node.Status == models.Registered {
			logrus.Infof("Node (%s) is already registered - skip", node.Mac)
			return
		}

		err = app.sql.GormDb.Model(&node).Update("status", models.Registered).Error
		if err != nil {
			logrus.Error(err)
			return
		}
		logrus.Infof("Update Node (MAC = %s) status to %s", node.Mac,
			models.Registered)

		// if valid, add into the unregistered Node db and send ack
		app.PublishAck(network.Name, payload.MAC)
	*/
}
