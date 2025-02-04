package hal

import (
	"bufio"
	"errors"
	"os"

	"github.com/aretas77/iot-controller/utils"
	"github.com/sirupsen/logrus"
)

const (
	// ActiveMode means radio chip is on. The chip can receive, transmit or listen.
	ActiveMode = "active"
	// ModemSleepMode means the CPU is operational and the clock is configurable.
	// WiFi/Bluetooth are disabled.
	ModemSleepMode = "modemSleep"
	// LightSleepMode means the CPU is paused. Any wake-up event will wake the device.
	LightSleepMode = "lightSleep"
	// DeepSleepMode ...
	DeepSleepMode = "deepSleep"
	// dataPath is used for storing sensor data.
	dataPath = "./cmd/data/"
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
	Protocol        string
	TemperatureLine int

	// Statistics simulation file related stuff.
	// Device should read from the file every 5 minutes as minimum, even if the
	// stats are not sent - this way the scanner will be incremented and
	// synchronized with values in iotctl service.
	StatisticsFileDesc *os.File
	StatisticsScanner  *bufio.Scanner
	StatisticsFrom     int
	StatisticsTo       int
	StatisticsCount    int
	StatisticsFileName string
}

func (e *ESP32) Initialize() error {
	f, err := os.Open(dataPath + e.StatisticsFileName)
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
	e.Protocol = "n"

	// Alright then, keep your file struct
	e.StatisticsFileDesc = f
	e.StatisticsScanner = bufio.NewScanner(e.StatisticsFileDesc)
	e.StatisticsCount = 0

	// Need to put the Scanner into correct place, if scanner starts from
	// the first line - we are already good to go.
	go func() {
		if e.StatisticsFrom != 0 {
			i := 1
			for e.StatisticsScanner.Scan() {
				e.StatisticsCount++
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
	}()

	logrus.Debug("initialized ESP32 HAL")
	return nil
}

// GetInterface will return the name of the underlying interface.
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

// GetPressureTemperature will read from the sensor
func (e *ESP32) GetPressureTemperature(sensor string) (float32, float32, float32) {
	var consumed float32

	switch sensor {
	case "bmp180":
		consumed += 0.007 // 7 micro Amperes for one reading, high res mode
	default:
	}

	// Simulate Temperature reading - read from a file
	e.StatisticsCount++
	e.StatisticsScanner.Scan()
	err, _, temp, press, _ := utils.SplitDataReadLine(e.StatisticsScanner.Text())
	if err != nil {
		logrus.Error(err)
	}

	return consumed, float32(temp), float32(press)
}

// SetPowerMode will set the operating power mode of the device. The device
// should use these modes to more accurately track of its power levels.
// TODO: track used power levels by power mode.
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

// GetSendConsumed ...
func (e *ESP32) GetSendConsumed() float32 {
	switch e.Protocol {
	case "n":
		return float32(e.Power.Tx80211n)
	case "g":
		return float32(e.Power.Tx80211g)
	case "b":
		return float32(e.Power.Tx80211b)
	default:
		return 0
	}
}

// PowerOff imitates device turn off - cleans data.
func (e *ESP32) PowerOff() {
	e.StatisticsFileDesc.Close()
}

func (e *ESP32) GetStatisticsInterval() (string, int, int) {
	return e.StatisticsFileName, e.StatisticsFrom, e.StatisticsTo
}

// GetStatisticsCurrentLine will return the line of the statistics file which
// was sent to the
func (e *ESP32) GetStatisticsCurrentLine() int {
	return e.StatisticsCount
}
