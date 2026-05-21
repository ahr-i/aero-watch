package db

import (
	"errors"
	"fmt"

	"github.com/ahr-i/aero-watch/drone-operation/setting"
)

var ErrDroneAlreadyExists = errors.New("drone model already exists")
var ErrInvalidStatus = errors.New("invalid drone status")

type Store interface {
	Init() error
	RegisterDroneModel(group string, code string) error
	ValidateDroneModel(group string, code string) error
	UpdateDroneStatus(group string, code string, status string) error
	CreateDriverInfo(content string) (DriverInfo, error)
	ListDriverInfos() ([]DriverInfo, error)
	UpdateDriverInfo(id int64, content string) error
	DeleteDriverInfo(id int64) error
	Close() error
}

type DriverInfo struct {
	ID      int64
	Content string
}

func NewStore() (Store, error) {
	switch setting.Setting.DB.Type {
	case "mysql":
		return NewMySQLStore(), nil
	default:
		return nil, fmt.Errorf("unsupported db type: %s", setting.Setting.DB.Type)
	}
}
