package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	db "github.com/aretas77/iot-controller/web/iotctl/database"
	models "github.com/aretas77/iot-controller/web/iotctl/database/models"
	mysql "github.com/aretas77/iot-controller/web/iotctl/database/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type NodeController struct {
	TableName string
	Database  *db.Database

	// Nodes will be saved at MySQL database so just keep a pointer into
	// MySql struct for easier access.
	sql *mysql.MySql
}

func (n *NodeController) Init() (err error) {
	if n.Database == nil {
		return errors.New("NodeController: Database is nil!")
	}

	if n.sql, err = n.Database.GetMySql(); err != nil {
		logrus.Error("NodeController: failed to get MySQL instance")
		return err
	}

	n.migrateNodeGorm()
	logrus.Debug("Initialized NodeController")
	return
}

func (n *NodeController) migrateNodeGorm() error {
	// Create a Node with additional settings
	settings := &models.NodeSettings{
		ReadInterval: 10,
		SendInterval: 10,
		NodeID:       2,
	}

	settings2 := &models.NodeSettings{
		ReadInterval: 10,
		SendInterval: 20,
		NodeID:       1,
	}

	if n.sql.GormDb.NewRecord(settings) {
		n.sql.GormDb.Create(&settings)
	}

	if n.sql.GormDb.NewRecord(settings2) {
		n.sql.GormDb.Create(&settings2)
	}

	node := &models.Node{
		Name:         "TestNode",
		Mac:          "AA:BB:CC:DD:EE",
		Location:     "Kaunas",
		IpAddress4:   "172.8.0.20",
		IpAddress6:   "NA",
		LastSentAck:  time.Now(),
		Status:       "acknowledged",
		SettingsID:   settings2.ID,
		NetworkRefer: 1,
	}

	node2 := &models.Node{
		Name:         "TestNode2",
		Mac:          "AA:BB:CC:DD:EF",
		Location:     "Kaunas",
		IpAddress4:   "172.8.0.21",
		IpAddress6:   "NA",
		LastSentAck:  time.Now(),
		Status:       "acknowledged",
		SettingsID:   settings.ID,
		NetworkRefer: 1,
	}

	if n.sql.GormDb.NewRecord(node) {
		n.sql.GormDb.Create(&node)
	}

	if n.sql.GormDb.NewRecord(node2) {
		n.sql.GormDb.Create(&node2)
	}

	return nil
}

func (n *NodeController) setupHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods",
		"POST, GET, OPTIONS, DELETE, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Access-Control-Allow-Origin")
}

// GetNode should return a Node by its ID.
func (n *NodeController) GetNode(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {
	n.setupHeader(&w)

	vars := mux.Vars(r)
	node := models.Node{}
	n.sql.GormDb.First(&node, vars["id"])

	if node.Mac != "" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(node)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// GetNodes should return all registered Nodes.
func (n *NodeController) GetNodes(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {
	n.setupHeader(&w)

	nodes := []models.Node{}
	n.sql.GormDb.Find(&nodes)

	// Gather NodeSettings and Network for all Nodes.
	for i, node := range nodes {
		settings := models.NodeSettings{}
		n.sql.GormDb.Where("id = ?", node.SettingsID).First(&settings)

		network := models.Network{}
		n.sql.GormDb.Where("id = ?", node.NetworkRefer).First(&network)

		nodes[i].Settings = settings
		nodes[i].Network = &network
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nodes)
}

// RegisterNode will parse `UnregisteredNode` from the request and check it against
// database to see if such a Node already exists and if it exists - it will
// register the `Node` and won't create an entry for `UnregisteredNode`.
//
// Otherwise, the `UnregisteredNode` will be added to the database and will point
// to the `Network` it was created and won't point to any of the `Node`s.
func (n *NodeController) RegisterNode(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {

	n.setupHeader(&w)
	decoder := json.NewDecoder(r.Body)

	var tmpNode models.UnregisteredNode
	if err := decoder.Decode(&tmpNode); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		// MAC address was not specified - don't continue.
		if strings.Compare("", tmpNode.Mac) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check if a Node with current MAC address exists.
		node := models.Node{}
		err := n.sql.GormDb.Where("mac = ?", tmpNode.Mac).Find(&node).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// The Node exists and its status is Registered - ignore this request.
		if node.Status != "" && node.Status == models.Registered {
			logrus.Infof("A Node(ID = %d) is already %s", node.ID, node.Status)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Node exists and is ready for registration - change Node status and
		// required fields.
		if node.Status == models.Acknowledged && node.NetworkRefer == tmpNode.NetworkRefer {
			logrus.Infof("A Node(ID = %d) is %s and is ready for registration",
				node.ID, node.Status)
			logrus.Infof("Update Node(ID = %d) status to %s", node.ID,
				models.Registered)

			err = n.sql.GormDb.Model(&node).Update("status", models.Registered).Error
			if err != nil {
				logrus.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Nodes status was updated successfully:
			// 'acknowledged' -> 'registered'.
			w.WriteHeader(http.StatusOK)
			return
		}

		// Node doesn't exist yet - add UnregisteredNode into a List for future
		// requests.
		err = n.sql.GormDb.Where(models.UnregisteredNode{Mac: tmpNode.Mac}).FirstOrCreate(&tmpNode).Error
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}

	// If there are no existing `Node` with the same MAC address as an
	// `UnregisteredNode` then a new `UnregisteredNode` entry is created in
	// a database and thus StatusCreated is returned.
	w.WriteHeader(http.StatusCreated)
}

// UnRegisterNode should remove the Node from our network.
func (n *NodeController) UnregisterNode(w http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {

}
