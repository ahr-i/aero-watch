package setting

var Setting settingStruct

type settingStruct struct {
	ServerPort           string `json:"server_port"`
	RTMPPort             string `json:"rtmp_port"`
	HLSRoot              string `json:"hls_root"`
	CaptureRoot          string `json:"capture_root"`
	HLSTimeSeconds       int    `json:"hls_time_seconds"`
	HLSListSize          int    `json:"hls_list_size"`
	StreamTimeoutSeconds int    `json:"stream_timeout_seconds"`
	DroneService         string `json:"drone_service"`
	DroneValidatePath    string `json:"drone_validate_path"`
	DroneValidateEnabled bool   `json:"drone_validate_enabled"`
	FFmpegPath           string `json:"ffmpeg_path"`
}
