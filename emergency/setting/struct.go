package setting

var Setting settingStruct

type settingStruct struct {
	ServerPort            string                `json:"serverPort"`
	CSVPath               string                `json:"csvPath"`
	DataFilePattern       string                `json:"dataFilePattern"`
	TablePrefix           string                `json:"tablePrefix"`
	NearestEmergencyLimit int                   `json:"nearestEmergencyLimit"`
	DB                    databaseSettingStruct `json:"db"`
	CSVSchema             map[string][]string   `json:"csvSchema"`
}

type databaseSettingStruct struct {
	HostEnv     string `json:"hostEnv"`
	PortEnv     string `json:"portEnv"`
	UserEnv     string `json:"userEnv"`
	PasswordEnv string `json:"passwordEnv"`
	SchemaEnv   string `json:"schemaEnv"`
}
