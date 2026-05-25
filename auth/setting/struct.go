package setting

var Setting settingStruct

type settingStruct struct {
	ServerPort string          `json:"serverPort"`
	DB         dbStruct        `json:"db"`
	JWT        jwtStruct       `json:"jwt"`
	Admin      adminStruct     `json:"admin"`
	Role       roleStruct      `json:"role"`
	UserTable  userTableStruct `json:"userTable"`
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

type jwtStruct struct {
	SecretEnv            string `json:"secretEnv"`
	AccessTokenExpireMin int    `json:"accessTokenExpireMin"`
	Secret               string `json:"-"`
}

type adminStruct struct {
	UserEnv     string `json:"userEnv"`
	PasswordEnv string `json:"passwordEnv"`
	User        string `json:"-"`
	Password    string `json:"-"`
}

type roleStruct struct {
	Admin      string `json:"admin"`
	Unverified string `json:"unverified"`
	Active     string `json:"active"`
}

type userTableStruct struct {
	Name               string `json:"name"`
	UsernameColumn     string `json:"usernameColumn"`
	PasswordHashColumn string `json:"passwordHashColumn"`
	RoleColumn         string `json:"roleColumn"`
	CreatedAtColumn    string `json:"createdAtColumn"`
}
