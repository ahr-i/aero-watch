package setting

var Setting settingStruct

type settingStruct struct {
	ServerPort string        `json:"serverPort"`
	Services   serviceStruct `json:"services"`
}

type serviceStruct struct {
	AIAnalysis     serviceConfig `json:"aiAnalysis"`
	Auth           serviceConfig `json:"auth"`
	DroneOperation serviceConfig `json:"droneOperation"`
	Emergency      serviceConfig `json:"emergency"`
	GPSTracking    serviceConfig `json:"gpsTracking"`
	Streaming      serviceConfig `json:"streaming"`
}

type serviceConfig struct {
	BaseURL string `json:"baseUrl"`
}
