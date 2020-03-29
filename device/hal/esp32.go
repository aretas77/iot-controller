package hal

import (
	"bufio"
	"errors"
	"os"

	"github.com/aretas77/iot-controller/utils"
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
	TemperatureLine int

	StatisticsFileDesc *os.File
	StatisticsScanner  *bufio.Scanner
	StatisticsFrom     int
	StatisticsTo       int
	StatisticsFileName string
}

func (e *ESP32) Initialize() error {
	f, err := os.Open("./cmd/data/" + e.StatisticsFileName)
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

	// Alright then, keep your file struct
	e.StatisticsFileDesc = f
	e.StatisticsScanner = bufio.NewScanner(e.StatisticsFileDesc)

	// Need to put the Scanner into correct place
	if e.StatisticsFrom != 0 {
		i := 1
		for e.StatisticsScanner.Scan() {
			if i >= e.StatisticsFrom && i < e.StatisticsTo {
				//logrus.Info(e.StatisticsScanner.Text())
				break
			}
			i++
		}

		if err := e.StatisticsScanner.Err(); err != nil {
			logrus.Error(err)
		}
	}

	logrus.Debugf("initialized ESP32 HAL")
	return nil
}

func (e *ESP32) GetInterface() string {
	return e.Interface
}

// GetTemperature will read from a sensor of a given name.
func (e *ESP32) GetTemperature(sensor string) (float32, float32) {
	var consumed float32

	switch sensor {
	case "bmp180":
		consumed += 0.007 // 7 micro Amperes for one reading, high res mode
	default:
	}

	// Simulate Temperature reading - read from a file
	e.StatisticsScanner.Scan()
	err, _, temp, _, _ := utils.SplitDataReadLine(e.StatisticsScanner.Text())
	if err != nil {
		logrus.Error(err)
	}

	return consumed, float32(temp)
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
