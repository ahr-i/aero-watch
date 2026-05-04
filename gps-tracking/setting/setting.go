package setting

import (
	"encoding/json"
	"os"

	"github.com/ahr-i/aero-watch/gps-tracking/utils/logging"
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

	setDefaultValues()

	return nil
}

func setDefaultValues() {
	if Setting.ServerReadHeaderTimeoutSec <= 0 {
		Setting.ServerReadHeaderTimeoutSec = 5
	}
	if Setting.ServerReadTimeoutSec <= 0 {
		Setting.ServerReadTimeoutSec = 10
	}
	if Setting.ServerWriteTimeoutSec <= 0 {
		Setting.ServerWriteTimeoutSec = 10
	}
	if Setting.ServerIdleTimeoutSec <= 0 {
		Setting.ServerIdleTimeoutSec = 60
	}
	if Setting.GPSAliveTimeoutSec <= 0 {
		Setting.GPSAliveTimeoutSec = 10
	}
	if Setting.GPSCleanupIntervalSec <= 0 {
		Setting.GPSCleanupIntervalSec = 1
	}
	if Setting.DroneValidateTimeoutSec <= 0 {
		Setting.DroneValidateTimeoutSec = 5
	}
}
