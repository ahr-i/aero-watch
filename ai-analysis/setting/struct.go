package setting

var Setting settingStruct

type settingStruct struct {
	ServerPort string       `json:"server_port"`
	OpenAI     openAIStruct `json:"openai"`
}

type openAIStruct struct {
	APIURL           string `json:"api_url"`
	Model            string `json:"model"`
	TimeoutSeconds   int    `json:"timeout_seconds"`
	APIKeyEnv        string `json:"api_key_env"`
	TextPrePromptEnv string `json:"text_pre_prompt_env"`
	TextPostPromptEnv string `json:"text_post_prompt_env"`
	ImagePrePromptEnv string `json:"image_pre_prompt_env"`
	ImagePostPromptEnv string `json:"image_post_prompt_env"`
	APIKey           string `json:"-"`
	TextPrePrompt    string `json:"-"`
	TextPostPrompt   string `json:"-"`
	ImagePrePrompt   string `json:"-"`
	ImagePostPrompt  string `json:"-"`
}
