package setting

var Setting settingStruct

type settingStruct struct {
	ServerPort string          `json:"server_port"`
	DB         dbStruct        `json:"db"`
	JWT        jwtStruct       `json:"jwt"`
	Admin      adminStruct     `json:"admin"`
	Role       roleStruct      `json:"role"`
	UserTable  userTableStruct `json:"user_table"`
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

type jwtStruct struct {
	SecretEnv            string `json:"secret_env"`
	AccessTokenExpireMin int    `json:"access_token_expire_min"`
	Secret               string `json:"-"`
}

type adminStruct struct {
	UserEnv     string `json:"user_env"`
	PasswordEnv string `json:"password_env"`
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
	UsernameColumn     string `json:"username_column"`
	PasswordHashColumn string `json:"password_hash_column"`
	RoleColumn         string `json:"role_column"`
}
