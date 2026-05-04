package setting

var Setting settingStruct

type settingStruct struct {
	ServerPort                 string `json:"server_port"`
	ServerReadHeaderTimeoutSec int    `json:"server_read_header_timeout_sec"`
	ServerReadTimeoutSec       int    `json:"server_read_timeout_sec"`
	ServerWriteTimeoutSec      int    `json:"server_write_timeout_sec"`
	ServerIdleTimeoutSec       int    `json:"server_idle_timeout_sec"`
	GPSAliveTimeoutSec         int    `json:"gps_alive_timeout_sec"`
	GPSCleanupIntervalSec      int    `json:"gps_cleanup_interval_sec"`
	DroneService               string `json:"drone_service"`
	DroneValidatePath          string `json:"drone_validate_path"`
	DroneValidateEnabled       bool   `json:"drone_validate_enabled"`
	DroneValidateTimeoutSec    int    `json:"drone_validate_timeout_sec"`
}
