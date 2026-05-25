package setting

import (
	"bufio"
	"os"
	"strings"
)

const envFilePath string = ".env"

func loadEnvFile() error {
	file, err := os.Open(envFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}
