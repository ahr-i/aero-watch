package setting

var Setting settingStruct

type settingStruct struct {
	ServerPort string `json:"serverPort"`
	DistPath   string `json:"distPath"`
}
