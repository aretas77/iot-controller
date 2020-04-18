package iotctl

import (
	"encoding/json"
	"time"

	"github.com/aretas77/iot-controller/types/devices"
	"github.com/aretas77/iot-controller/types/mqtt"
	"github.com/aretas77/iot-controller/utils"
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
		if !gorm.IsRecordNotFoundError(err) {
			logrus.Error(err)
		}

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

	settings := models.NodeSettings{ReadInterval: tmpNode.InitialReadInterval}
	if err = app.sql.GormDb.Create(&settings).Error; err != nil {
		panic(err)
	}

	node := models.Node{
		Name:                payload.Name,
		Mac:                 payload.MAC,
		Location:            "Kaunas",
		IpAddress4:          payload.IpAddress4,
		LastSentAck:         payload.Sent,
		LastReceivedMessage: time.Now(),
		Status:              models.Registered,
		SettingsID:          settings.ID,
		NetworkRefer:        tmpNode.NetworkRefer,
		AddedUsername:       tmpNode.AddedUsername,
	}
	app.sql.GormDb.Create(&node)
	app.sql.GormDb.Model(&settings).Update("node_id", node.ID)

	// Delete permanently
	if err = app.sql.GormDb.Unscoped().Delete(&tmpNode).Error; err != nil {
		panic(err)
	}

	logrus.Infof("Created a new Node (ID = %d, MAC = %s)", node.ID, node.Mac)

	// We can publish ack when we have done the following:
	//	1. An `UnregisteredNode` existed for this particular device node.
	//	2. We have created a new `Node` for this device with its status
	//	   as 'Registered'.
	app.PublishAck(network.Name, payload.MAC, node.Location, tmpNode.InitialReadInterval)
	return
}

// OnMessageStats will handle a received sensor data from the device. When a
// message is received, it should check whether a node for this sensor data
// exists and if it does - add to the database and forward the sensor data to
// the IoT Hades service for processing.
// Topic: node/+/+/stats
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

	// parse MAC address from the topic
	_, _, mac, _ := utils.SplitTopic4(msg.Topic())

	// Check if such Node exists, and if it exists - update the statistics.
	node, err := app.sql.CheckNodeExists(mac)
	if err != nil {
		logrus.Error(err)
		return
	}

	entry := models.NodeStatisticsEntry{
		CPULoad:           payload.CPULoad,
		Pressure:          payload.Pressure,
		Temperature:       payload.Temperature,
		TempReadTime:      payload.TempReadTime,
		Consumed:          payload.Consumed,
		NodeRefer:         node.Mac,
		BatteryMah:        payload.BatteryLeft,
		BatteryPercentage: payload.BatteryPercentage,
		DataStatsLine:     payload.StatisticsCount,
		SendTimes:         payload.SendTimes,
		SensorReadTimes:   payload.SensorReadTimes,
	}

	if err := app.sql.GormDb.Create(&entry).Error; err != nil {
		logrus.Error(err)
		return
	}

	// update node values
	if err := app.sql.GormDb.Model(&node).Update(models.Node{
		BatteryMah:        payload.BatteryLeft,
		BatteryPercentage: payload.BatteryPercentage,
	}).Error; err != nil {
		logrus.Error(err)
		return
	}

	return
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

	node, err := app.sql.CheckNodeExists(payload.Mac)
	if err != nil {
		logrus.Error(err)
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

	// Get the `NodeSettings` of the particular `Node`. Such entry should exist.
	nodeSettings := models.NodeSettings{}
	err = app.sql.GormDb.Where("node_id = ?", node.ID).Find(&nodeSettings).Error
	if err != nil {
		logrus.Error(err)
		return
	}

	nodeSettings.DataFileName = payload.DataFileInfo.Filename
	nodeSettings.DataLineFrom = payload.DataFileInfo.DataLineFrom
	nodeSettings.DataLineTo = payload.DataFileInfo.DataLineTo

	// We arrive here if:
	//	1. `Node` with MAC exists in the DB.
	//	2. `Node` status is 'Registered' both in DB and on the device.
	//	3. `Node` networks are the same.
	// So now we can update the Nodes fields.
	node.LastReceivedMessage = time.Now()
	node.BatteryMah = payload.BatteryMah
	node.BatteryPercentage = payload.BatteryPercentage

	// Assume that the first received battery size is the full battery size
	if node.BatteryMahTotal == 0 && node.BatteryPercentage > 95 {
		node.BatteryMahTotal = payload.BatteryMah
	}

	// Need to update IP address in case something has changed.
	node.IpAddress4 = payload.IpAddress4

	// Update `Node` in the Database with new values.
	if err := app.sql.GormDb.Save(&node).Error; err != nil {
		logrus.Error(err)
		return
	}

	// Update `NodeSettings` in the Database with new values.
	if err := app.sql.GormDb.Save(&nodeSettings).Error; err != nil {
		logrus.Error(err)
		return
	}

	return
}

// OnMessageEvent will handle the events receveived by Hades daemon.
func (app *Iotctl) OnMessageEvent(client MQTT.Client, msg MQTT.Message) {
	logrus.Infof("plain message got on: %s", msg.Topic())
	payload := mqtt.MessageEventSent{}

	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"topic": msg.Topic,
			"msg":   msg.Payload,
		}).Error("failed to unmarshal greeting message")
		return
	}

	_, _, mac, _ := utils.SplitTopic4(msg.Topic())

	event := models.Event{
		Mac:      mac,
		Model:    payload.Model,
		TimeSent: payload.TimeSent,
	}

	if err := app.sql.GormDb.Create(&event).Error; err != nil {
		logrus.Error(err)
		return
	}

	return
}
