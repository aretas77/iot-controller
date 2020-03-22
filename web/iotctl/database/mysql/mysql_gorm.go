package mysql

import (
	models "github.com/aretas77/iot-controller/web/iotctl/database/models"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func (m *MySql) ConnectGorm() (err error) {
	logrus.Debug("Setting up MySQL database using GORM")

	m.GormDb, err = gorm.Open("mysql", m.Server)
	if err != nil {
		logrus.Error(ErrOpenFailed)
		panic(err.Error())
	}

	m.GormDb.LogMode(true)

	// Get the generic database object sql.DB to use its functions
	m.Db = m.GormDb.DB()

	logrus.Infof("Connected to MySQL at %s", m.Server)
	return
}

// InitializeMigrationGorm will create a database structure so it would be
// possible to manipulate data with it.
func (m *MySql) InitializeMigrationGorm() {
	m.GormDb.Model(&models.Node{}).RemoveForeignKey("settings_id", "node_settings(id)")
	m.GormDb.Model(&models.Node{}).RemoveForeignKey("network_refer", "networks(id)")
	m.GormDb.Model(&models.UnregisteredNode{}).RemoveForeignKey("network_refer", "networks(id)")
	m.GormDb.Model(&models.Network{}).RemoveForeignKey("user_refer", "users(id)")
	m.GormDb.Model(&models.NodeStatisticsEntry{}).RemoveForeignKey("node_refer", "nodes(id)")

	m.GormDb.DropTableIfExists(&models.Network{}, &models.User{},
		&models.Node{}, &models.NodeSettings{}, &models.UnregisteredNode{},
		&models.NodeStatisticsEntry{})

	m.GormDb.CreateTable(&models.User{}, &models.Node{}, &models.NodeSettings{},
		&models.UnregisteredNode{}, &models.Network{}, &models.NodeStatisticsEntry{})

	m.GormDb.Model(&models.Node{}).AddForeignKey("settings_id",
		"node_settings(id)", "RESTRICT", "RESTRICT")
	m.GormDb.Model(&models.Node{}).AddForeignKey("network_refer",
		"networks(id)", "RESTRICT", "RESTRICT")
	m.GormDb.Model(&models.Network{}).AddForeignKey("user_refer",
		"users(id)", "RESTRICT", "RESTRICT")
	m.GormDb.Model(&models.UnregisteredNode{}).AddForeignKey("network_refer",
		"networks(id)", "RESTRICT", "RESTRICT")
	m.GormDb.Model(&models.NodeStatisticsEntry{}).AddForeignKey("node_refer",
		"nodes(id)", "RESTRICT", "RESTRICT")
	//m.GormDb.Model(&models.UnregisteredNode{}).AddForeignKey("mac",
	//"nodes(mac)", "RESTRICT", "RESTRICT")

}

// CheckAuth should check whether given credentials are valid and if valid,
// return the User.
func (m *MySql) CheckUserExists(creds *models.Credentials) (*models.User, error) {
	var user models.User

	logrus.Debugf("Authenticating user (email = %s)", creds.Email)

	query := "email = ? AND password = ?"
	err := m.GormDb.Where(query, creds.Email, creds.Password).Find(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, models.ErrUserNotFound
		}
		return nil, models.ErrUserUnauthorized
	}

	logrus.Debugf("User (name = %s, email = %s) authenticated", user.Username,
		user.Email)
	return &user, nil
}

// CheckNodeExists will check whether a `Node` exists with a given MAC address
// and return a Node object if it does exist.
func (m *MySql) CheckNodeExists(mac string) (*models.Node, error) {
	var node models.Node

	err := m.GormDb.Where("mac = ?", mac).Find(&node).Error
	if err != nil {
		return nil, models.ErrNodeNotFound
	}

	return &node, nil
}
