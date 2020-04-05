package devices

// System is the configuration of the whole system and is periodically sent
// to the server.
type System struct {
	Mac                  string                `json:"mac"`
	Name                 string                `json:"name"`
	Location             string                `json:"location"`
	Network              string                `json:"network"`
	IpAddress4           string                `json:"ip_address4"`
	Status               string                `json:"status"`
	BatteryPercentage    float32               `json:"battery_left_per"`
	BatteryMah           float32               `json:"battery_left_mah"`
	CurrentBatteryMah    float32               `json:"cur_battery_mah"`
	FreeNetworkScheduler *FreeNetworkScheduler `json:"free_net_scheduler,omitempty"`
	DataFileInfo         DataFileInfo          `json:"data_file_info"`
}

type DataFileInfo struct {
	Filename     string `json:"filename"`
	DataLineFrom int    `json:"from"`
	DataLineTo   int    `json:"to"`
	DataCount    int    `json:"count"`
}

// FreeNetworkScheduler ...
type FreeNetworkScheduler struct {
	Enabled   bool                      `json:"enabled"`
	Monday    FreeNetworkSchedulerRange `json:"monday"`
	Tuesday   FreeNetworkSchedulerRange `json:"tuesday"`
	Wednesday FreeNetworkSchedulerRange `json:"wednesday"`
	Thursday  FreeNetworkSchedulerRange `json:"thursday"`
	Friday    FreeNetworkSchedulerRange `json:"friday"`
	Saturday  FreeNetworkSchedulerRange `json:"saturday"`
	Sunday    FreeNetworkSchedulerRange `json:"sunday"`
}

type FreeNetworkSchedulerRange struct {
	Disabled bool   `json:"disabled"`
	From     string `json:"from"`
	To       string `json:"to"`
}
