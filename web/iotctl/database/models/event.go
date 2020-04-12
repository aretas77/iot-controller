package models

// Event represents a singular event of a model being sent by Hades to the
// `NodeDevice`.
type Event struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Model    string `json:"model"`
	TimeSent string `json:"time_sent"`
	Mac      string `json:"mac" gorm:"not null`
}
