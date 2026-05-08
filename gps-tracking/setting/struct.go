package setting

var Setting settingStruct

type settingStruct struct {
	ServerPort                 string `json:"serverPort"`
	ServerReadHeaderTimeoutSec int    `json:"serverReadHeaderTimeoutSec"`
	ServerReadTimeoutSec       int    `json:"serverReadTimeoutSec"`
	ServerWriteTimeoutSec      int    `json:"serverWriteTimeoutSec"`
	ServerIdleTimeoutSec       int    `json:"serverIdleTimeoutSec"`
	GPSAliveTimeoutSec         int    `json:"gpsAliveTimeoutSec"`
	GPSCleanupIntervalSec      int    `json:"gpsCleanupIntervalSec"`
	DroneService               string `json:"droneService"`
	DroneValidatePath          string `json:"droneValidatePath"`
	DroneValidateEnabled       bool   `json:"droneValidateEnabled"`
	DroneValidateTimeoutSec    int    `json:"droneValidateTimeoutSec"`
}
