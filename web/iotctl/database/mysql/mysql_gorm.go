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

	m.GormDb.DropTableIfExists(&models.Network{}, &models.User{},
		&models.Node{}, &models.NodeSettings{}, &models.UnregisteredNode{})

	m.GormDb.CreateTable(&models.User{}, &models.Node{}, &models.NodeSettings{},
		&models.UnregisteredNode{}, &models.Network{})

	m.GormDb.Model(&models.Node{}).AddForeignKey("settings_id",
		"node_settings(id)", "RESTRICT", "RESTRICT")
	m.GormDb.Model(&models.Node{}).AddForeignKey("network_refer",
		"networks(id)", "RESTRICT", "RESTRICT")
	m.GormDb.Model(&models.Network{}).AddForeignKey("user_refer",
		"users(id)", "RESTRICT", "RESTRICT")
	m.GormDb.Model(&models.UnregisteredNode{}).AddForeignKey("network_refer",
		"networks(id)", "RESTRICT", "RESTRICT")
	//m.GormDb.Model(&models.UnregisteredNode{}).AddForeignKey("mac",
	//"nodes(mac)", "RESTRICT", "RESTRICT")

}

// CheckAuth should check whether given credentials are valid and if valid,
// return the User.
func (m *MySql) CheckAuth(creds *models.Credentials) (*models.User, error) {
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
