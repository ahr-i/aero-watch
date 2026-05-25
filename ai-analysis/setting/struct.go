package setting

var Setting settingStruct

type settingStruct struct {
	ServerPort string       `json:"serverPort"`
	HistoryDir string       `json:"historyDir"`
	OpenAI     openAIStruct `json:"openai"`
}

type openAIStruct struct {
	APIURL              string `json:"apiUrl"`
	Model               string `json:"model"`
	TimeoutSeconds      int    `json:"timeoutSeconds"`
	APIKeyEnv           string `json:"apiKeyEnv"`
	TextPrePromptEnv    string `json:"textPrePromptEnv"`
	TextPostPromptEnv   string `json:"textPostPromptEnv"`
	ImagePrePromptEnv   string `json:"imagePrePromptEnv"`
	ImagePostPromptEnv  string `json:"imagePostPromptEnv"`
	APIKey              string `json:"-"`
	TextPrePrompt       string `json:"-"`
	TextPostPrompt      string `json:"-"`
	ImagePrePrompt      string `json:"-"`
	ImagePostPrompt     string `json:"-"`
}
