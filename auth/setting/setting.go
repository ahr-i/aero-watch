package setting

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ahr-i/aero-watch/auth/utils/logging"
)

const settingFilePath string = "./setting/setting.json"
const envFilePath string = "./.env"

func Init() {
	err := readSettingFile()
	if err != nil {
		logging.Error(err)

		os.Exit(1)
	}

	err = readEnvFile()
	if err != nil {
		logging.Error(err)

		os.Exit(1)
	}

	err = applyEnv()
	if err != nil {
		logging.Error(err)

		os.Exit(1)
	}

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

func readEnvFile() error {
	file, err := os.Open(envFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			return fmt.Errorf("invalid env line: %s", line)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		if key == "" {
			return fmt.Errorf("empty env key")
		}

		os.Setenv(key, value)
	}

	return scanner.Err()
}

func applyEnv() error {
	dbUser, err := requiredEnv(Setting.DB.UserEnv)
	if err != nil {
		return err
	}

	dbPassword, err := requiredEnv(Setting.DB.PasswordEnv)
	if err != nil {
		return err
	}

	jwtSecret, err := requiredEnv(Setting.JWT.SecretEnv)
	if err != nil {
		return err
	}

	adminUser, err := requiredEnv(Setting.Admin.UserEnv)
	if err != nil {
		return err
	}

	adminPassword, err := requiredEnv(Setting.Admin.PasswordEnv)
	if err != nil {
		return err
	}

	Setting.DB.User = dbUser
	Setting.DB.Password = dbPassword
	Setting.JWT.Secret = jwtSecret
	Setting.Admin.User = adminUser
	Setting.Admin.Password = adminPassword

	return nil
}

func requiredEnv(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("empty env setting")
	}

	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("missing required env: %s", key)
	}

	return value, nil
}
