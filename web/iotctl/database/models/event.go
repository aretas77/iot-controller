package models

// Event ...
type Event struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Model    string `json:"model"`
	TimeSent string `json:"time_sent"`
	Mac      string `json:"mac"`
}
