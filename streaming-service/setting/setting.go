package setting

import (
	"encoding/json"
	"os"

	"github.com/ahr-i/aero-watch/utils/logging"
)

const settingFilePath string = "./setting/setting.json"

func Init() {
	err := readSettingFile()
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
