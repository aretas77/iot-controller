package hal

import (
	"bufio"
	"errors"
	"os"

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

// ESP32 is a HAL implementation of a real life esp32 device and is used to
// simulate it.
// TODO: make PowerConsumption dependent on Mode.
type ESP32 struct {
	Interface       string // Name of this struct
	Power           *PowerConsumption
	Mode            string
	Statistics      string
	TemperatureLine int

	StatisticsFileDesc *os.File
	StatisticsScanner  *bufio.Reader
}

func (e *ESP32) Initialize(statFile string) error {
	f, err := os.Open("./cmd/data/" + statFile)
	if err != nil {
		panic(err)
	}

	e.Power = &PowerConsumption{
		// All values in mA, taken from ESP32 datasheet.
		// https://www.espressif.com/sites/default/files/documentation/esp32_datasheet_en.pdf, page 42.
		Tx80211b:   240,
		Tx80211g:   190,
		Tx80211n:   180,
		Rx80211bgn: 100,
		TxBt:       130,
		RxBt:       100,
	}
	e.Mode = ActiveMode
	e.Interface = "esp32"
	e.Statistics = statFile

	// Alright then, keep your file struct
	e.StatisticsFileDesc = f
	e.StatisticsScanner = bufio.NewReader(e.StatisticsFileDesc)

	logrus.Debugf("initialized ESP32 HAL")
	return nil
}

func (e *ESP32) GetDeviceName() string {
	return ""
}

func (e *ESP32) GetInterface() string {
	return e.Interface
}

// GetTemperature will read from a sensor of a given name.
func (e *ESP32) GetTemperature(sensor string) (float32, float32) {
	var consumed, temperature float32

	switch sensor {
	case "bmp180":
		consumed += 0.007 // 7 micro Amperes for one reading, high res mode
	default:
	}

	if e.StatisticsScanner == nil {
		logrus.Info("nil")
	}

	line, err := e.StatisticsScanner.ReadBytes('\n')
	if err != nil {
		panic(err)
	}
	logrus.Info(string(line))

	temperature = 10

	return consumed, temperature
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

func (e *ESP32) PowerOff() {
	e.StatisticsFileDesc.Close()
}
