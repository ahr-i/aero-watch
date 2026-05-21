package setting

var Setting settingStruct

type settingStruct struct {
	ServerPort string           `json:"serverPort"`
	DB         dbStruct         `json:"db"`
	DroneTable droneTableStruct `json:"droneTable"`
	Status     statusStruct     `json:"status"`
}

type dbStruct struct {
	Type        string `json:"type"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Name        string `json:"name"`
	UserEnv     string `json:"userEnv"`
	PasswordEnv string `json:"passwordEnv"`
	User        string `json:"-"`
	Password    string `json:"-"`
}

type droneTableStruct struct {
	Name         string `json:"name"`
	GroupColumn  string `json:"groupColumn"`
	CodeColumn   string `json:"codeColumn"`
	StatusColumn string `json:"statusColumn"`
}

type statusStruct struct {
	Active   string `json:"active"`
	Inactive string `json:"inactive"`
}
