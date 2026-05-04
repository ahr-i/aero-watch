package setting

var Setting settingStruct

type settingStruct struct {
	ServerPort string           `json:"server_port"`
	DB         dbStruct         `json:"db"`
	DroneTable droneTableStruct `json:"drone_table"`
	Status     statusStruct     `json:"status"`
}

type dbStruct struct {
	Type        string `json:"type"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Name        string `json:"name"`
	UserEnv     string `json:"user_env"`
	PasswordEnv string `json:"password_env"`
	User        string `json:"-"`
	Password    string `json:"-"`
}

type droneTableStruct struct {
	Name         string `json:"name"`
	GroupColumn  string `json:"group_column"`
	CodeColumn   string `json:"code_column"`
	StatusColumn string `json:"status_column"`
}

type statusStruct struct {
	Active   string `json:"active"`
	Inactive string `json:"inactive"`
}
