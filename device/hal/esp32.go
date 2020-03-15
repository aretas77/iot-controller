package hal

import (
	"errors"

	"github.com/sirupsen/logrus"
)

const (
	// Radio chip is on. The chip can receive, transmit or listen.
	ActiveMode = "active"

	// The CPU is operational and the clock is configurable. WiFi/Bluetooth
	// are disabled.
	ModemSleepMode = "modemSleep"

	// The CPU is paused. Any wake-up event will wake the device.
	LightSleepMode = "lightSleep"

	//
	DeepSleepMode = "deepSleep"
)

var (
	ErrModeInvalid = errors.New("mode doesn't exist")
)

type ESP32 struct {
	Interface string // Name of this struct

	// TODO: make PowerConsumption dependent on Mode.
	Power   *PowerConsumption
	Mode    string
	Battery int
}

func (e *ESP32) Initialize() error {
	e = &ESP32{
		Power: &PowerConsumption{
			// All values in mA, taken from ESP32 datasheet.
			// https://www.espressif.com/sites/default/files/documentation/esp32_datasheet_en.pdf, page 42.
			Tx80211b:   240,
			Tx80211g:   190,
			Tx80211n:   180,
			Rx80211bgn: 100,
			TxBt:       130,
			RxBt:       100,
		},
		Mode:      ActiveMode,
		Interface: "esp32",
		// Battery at 100%.
		Battery: 100,
	}

	logrus.Debugf("initialized ESP32 HAL")
	return nil
}

func (e *ESP32) GetDeviceName() string {
	return ""
}

func (e *ESP32) GetInterface() string {
	return e.Interface
}

func (e *ESP32) GetTemperature() (float32, float32) {
	return 0, 0
}

func (e *ESP32) SetPowerMode(mode string) error {
	switch mode {
	case ActiveMode:
		e.Mode = ActiveMode
	case ModemSleepMode:
		e.Mode = ModemSleepMode
	case LightSleepMode:
		e.Mode = LightSleepMode
	case DeepSleepMode:
		e.Mode = DeepSleepMode
	default:
		return ErrModeInvalid
	}

	return nil
}

func (e *ESP32) GetPowerMode() string {
	return e.Mode
}

func (e *ESP32) GetBatteryPercentage() float32 {
	return 0
}
