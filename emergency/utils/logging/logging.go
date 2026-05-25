package logging

import (
	"fmt"
	"log"

	ipfslog "github.com/ipfs/go-log/v2"
)

var (
	systemLogger = ipfslog.Logger("system")
	debugLogger  = ipfslog.Logger("debug")
)

func Init() {
	level, err := ipfslog.LevelFromString("info")
	if err != nil {
		log.Println(err)
		return
	}

	ipfslog.SetAllLoggers(level)
	Info("Successfully initialized logging.")
}

func Info(msg interface{}) {
	systemLogger.Info(formatMessage(msg))
}

func Warn(msg interface{}) {
	systemLogger.Warn(formatMessage(msg))
}

func Error(msg interface{}) {
	systemLogger.Error(formatMessage(msg))
}

func Debug(msg interface{}) {
	debugLogger.Debug(formatMessage(msg))
}

func formatMessage(msg interface{}) string {
	switch v := msg.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return fmt.Sprintf("%v", v)
	}
}
