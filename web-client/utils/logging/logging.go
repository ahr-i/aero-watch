package logging

import (
	"fmt"
	"log"
)

func Init() {
	Info("Successfully initialized logging.")
}

func Info(msg interface{}) {
	log.Printf("[INFO] %s", formatMessage(msg))
}

func Warn(msg interface{}) {
	log.Printf("[WARN] %s", formatMessage(msg))
}

func Error(msg interface{}) {
	log.Printf("[ERROR] %s", formatMessage(msg))
}

func Debug(msg interface{}) {
	log.Printf("[DEBUG] %s", formatMessage(msg))
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
