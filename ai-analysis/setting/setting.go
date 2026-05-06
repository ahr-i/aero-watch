package setting

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/ahr-i/aero-watch/ai-analysis/utils/logging"
)

const settingFilePath string = "./setting/setting.json"
const envFilePath string = "./.env"

func Init() {
	err := readSettingFile()
	if err != nil {
		logging.Error(err)

		os.Exit(1)
	}

	env, err := readEnvFile()
	if err != nil && !os.IsNotExist(err) {
		logging.Error(err)

		os.Exit(1)
	}

	readOpenAIEnv(env)
	logging.Info("Successfully finished initializing setting.")
}

func readSettingFile() error {
	file, err := os.ReadFile(settingFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &Setting)
	if err != nil {
		return err
	}

	return nil
}

func readEnvFile() (map[string]string, error) {
	file, err := os.ReadFile(envFilePath)
	if err != nil {
		return nil, err
	}

	env := make(map[string]string)
	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		if key == "" {
			continue
		}

		env[key] = value
	}

	return env, nil
}

func readOpenAIEnv(env map[string]string) {
	Setting.OpenAI.APIKey = readEnvValue(env, Setting.OpenAI.APIKeyEnv)
	Setting.OpenAI.TextPrePrompt = readEnvValue(env, Setting.OpenAI.TextPrePromptEnv)
	Setting.OpenAI.TextPostPrompt = readEnvValue(env, Setting.OpenAI.TextPostPromptEnv)
	Setting.OpenAI.ImagePrePrompt = readEnvValue(env, Setting.OpenAI.ImagePrePromptEnv)
	Setting.OpenAI.ImagePostPrompt = readEnvValue(env, Setting.OpenAI.ImagePostPromptEnv)
}

func readEnvValue(env map[string]string, key string) string {
	if key == "" {
		return ""
	}

	return env[key]
}
