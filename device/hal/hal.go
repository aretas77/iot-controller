package hal

// PowerConsumption structure defines what can consume power and will
// specify how much.
// All values are in mA (mili-Amperes).
type PowerConsumption struct {
	Tx80211b   int // Transmit 802.11b
	Tx80211g   int // Transmit 802.11g
	Tx80211n   int // Transmit 802.11n
	Rx80211bgn int // Receive 802.11b/g/n
	TxBt       int // Transmit BT/BLE
	RxBt       int // Receive BT/BLE
}

// HAL is short for Hardware Abstraction Layer and its used as a reference
// to various devices with `built-in` sensors.
// TODO: `built-in` sensors should become as modules which could be dynamically
// attached to a device.
type HAL interface {

	// Initialize will setup and initialize the simulated device.
	Initialize() error

	// GetInterface will return the name of the interface e.g. esp32.
	GetInterface() string

	// GetTemperature will return a temperature value and how much energy was
	// used by using the sensor.
	//
	// @sensor - name of the sensor from which we will read.
	GetTemperature(sensor string) (float32, float32)

	// GetPressureTemperature will return a pressure and a temperature value
	// and how much energy was consumed by using the sensor.
	//
	// @sensor - name of the sensor from which we will read.
	GetPressureTemperature(sensor string) (float32, float32, float32)

	// GetSendConsumed will return how much power was consumed by sending a
	// packet. It is calculated by used protocol (as in `PowerConsumption`).
	GetSendConsumed() float32

	// SetPowerMode will set the power mode of the simulation device which
	// will adjust TX, RX and other peripheral power consumption.
	//
	// @mode - the mode in which the device operates in.
	SetPowerMode(mode string) error

	// GetPowerMode will return the devices operating mode.
	GetPowerMode() string

	// PowerOff should clean up.
	PowerOff()

	// GetStatisticsInterval will return the filename and the interval on which
	// the device operates by reading the data.
	GetStatisticsInterval() (string, int, int)

	// GetStatisticsCurrentLine ...
	GetStatisticsCurrentLine() int
}
