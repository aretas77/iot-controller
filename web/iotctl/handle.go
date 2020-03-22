package iotctl

import (
	"encoding/json"
	"time"

	"github.com/aretas77/iot-controller/types/devices"
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

// OnMessageGreeting ...
func (app *Iotctl) OnMessageGreeting(client MQTT.Client, msg MQTT.Message) {
	logrus.Infof("plain got message on: %s", msg.Topic())
	payload := mqtt.MessageGreeting{}

	// check if a valid network ID is provided
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"topic": msg.Topic,
			"msg":   msg.Payload,
		}).Error("failed to unmarshal greeting message")
		return
	}

	// check if a an `UnregisteredNode` is created with this MAC and depending
	// on whether it exists or not, branch out into two stages:
	//	1. `UnregisteredNode` doesn't exist - this means that the user hasn't
	//		requested any device with such MAC so we should ignore this devices
	//		greeting.
	//	2. `UnregisteredNode` with MAC exists - we know that the user has
	//		requested for this Node - we can send the ACK to the device and
	//		create a `Node` with fields: MAC.

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
		Name:                "AckNode",
		Mac:                 payload.MAC,
		Location:            "",
		IpAddress4:          payload.IpAddress4,
		LastSentAck:         payload.Sent,
		LastReceivedMessage: time.Now(),
		Status:              models.Registered,
		SettingsID:          settings.ID,
		NetworkRefer:        tmpNode.NetworkRefer,
		AddedUsername:       tmpNode.AddedUsername,
	}
	app.sql.GormDb.Create(&node)

	app.sql.GormDb.Model(&settings).Update("node_refer", node.ID)
	app.sql.GormDb.Delete(&tmpNode)

	logrus.Infof("Created a new Node (ID = %d, MAC = %s)", node.ID, node.Mac)

	// We can publish ack when we have done the following:
	//	1. An `UnregisteredNode` existed for this particular device node.
	//	2. We have created a new `Node` for this device with its status
	//	   as 'Registered'.
	app.PublishAck(network.Name, payload.MAC, node.Location)
	return
}

// OnMessageStats ...
func (app *Iotctl) OnMessageStats(client MQTT.Client, msg MQTT.Message) {
	logrus.Infof("plain got message on: %s", msg.Topic())
	payload := mqtt.MessageStats{}

	// check if a valid network ID is provided
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"topic": msg.Topic,
			"msg":   msg.Payload,
		}).Error("failed to unmarshal stats message")
		return
	}

	// Check if such Node exists, and if it exists - update the statistics.

}

// OnMessageSystem will handle the received messages about a currently running
// device status, e.g. its battery level, what HAL is used and etc.
func (app *Iotctl) OnMessageSystem(client MQTT.Client, msg MQTT.Message) {
	logrus.Infof("plain got message on: %s", msg.Topic())
	logrus.Debugf("received system information = %s", msg.Payload())
	payload := devices.System{}

	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"topic": msg.Topic,
			"msg":   msg.Payload,
		}).Error("failed to unmarshal stats message")
		return
	}

	node := models.Node{}
	err := app.sql.GormDb.Where("mac = ?", payload.Mac).Find(&node).Error
	if gorm.IsRecordNotFoundError(err) || err != nil {
		logrus.Error(err)
		// No such Node exists - skip it
		return
	}

	// We need to make sure that the Node is registered - greeting and ack done.
	if node.Status != models.Registered || payload.Status != models.Registered {
		logrus.Debugf("node (%s) is not registered - skip", payload.Mac)
		return
	}

	network := models.Network{}
	err = app.sql.GormDb.Where("id = ?", node.NetworkRefer).Find(&network).Error
	if err != nil {
		logrus.Error(err)
		return
	}

	// Are networks the same?
	if network.Name != payload.Network {
		logrus.Debugf("node (%s) network is wrong - %s != %s", node.Mac,
			network.Name, payload.Network)
		return
	}

	// We arrive here if:
	//	1. `Node` with MAC exists in the DB.
	//	2. `Node` status is 'Registered' both in DB and on the device.
	//	3. `Node` networks are the same.
	// So now we can update the Nodes fields.
	node.LastReceivedMessage = time.Now()
	node.BatteryMah = payload.BatteryMah
	node.BatteryPercentage = payload.BatteryPercentage

	// Need to update IP address in case something has changed.
	node.IpAddress4 = payload.IpAddress4

	// Update `Node` in the Database with new values.
	if err := app.sql.GormDb.Save(&node).Error; err != nil {
		logrus.Error(err)
		return
	}

	return
}
