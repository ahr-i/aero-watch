package setting

var Setting settingStruct

type settingStruct struct {
	ServerPort            string                      `json:"serverPort"`
	DB                    dbStruct                    `json:"db"`
	DroneTable            droneTableStruct            `json:"droneTable"`
	DriverInfoTable       driverInfoTableStruct       `json:"driverInfoTable"`
	DroneDriverMatchTable droneDriverMatchTableStruct `json:"droneDriverMatchTable"`
	Status                statusStruct                `json:"status"`
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

type driverInfoTableStruct struct {
	Name          string `json:"name"`
	IDColumn      string `json:"idColumn"`
	ContentColumn string `json:"contentColumn"`
}

type droneDriverMatchTableStruct struct {
	Name           string `json:"name"`
	DriverIDColumn string `json:"driverIdColumn"`
	GroupColumn    string `json:"groupColumn"`
	CodeColumn     string `json:"codeColumn"`
}

type statusStruct struct {
	Active   string `json:"active"`
	Inactive string `json:"inactive"`
}
