package setting

var Setting settingStruct

type settingStruct struct {
	ServerPort            string `json:"serverPort"`
	RTMPPort              string `json:"rtmpPort"`
	HLSRoot               string `json:"hlsRoot"`
	CaptureRoot           string `json:"captureRoot"`
	HLSTimeSeconds        int    `json:"hlsTimeSeconds"`
	HLSListSize           int    `json:"hlsListSize"`
	StreamTimeoutSeconds  int    `json:"streamTimeoutSeconds"`
	DroneOperationService string `json:"droneOperationService"`
	DroneValidatePath     string `json:"droneValidatePath"`
	DroneValidateEnabled  bool   `json:"droneValidateEnabled"`
	FFmpegPath            string `json:"ffmpegPath"`
}
