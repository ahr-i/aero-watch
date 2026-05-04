package db

import (
	"errors"
	"fmt"

	"github.com/ahr-i/aero-watch/drone/setting"
)

var ErrDroneAlreadyExists = errors.New("drone model already exists")
var ErrInvalidStatus = errors.New("invalid drone status")

type Store interface {
	Init() error
	RegisterDroneModel(group string, code string) error
	ValidateDroneModel(group string, code string) error
	UpdateDroneStatus(group string, code string, status string) error
	Close() error
}

func NewStore() (Store, error) {
	switch setting.Setting.DB.Type {
	case "mysql":
		return NewMySQLStore(), nil
	default:
		return nil, fmt.Errorf("unsupported db type: %s", setting.Setting.DB.Type)
	}
}
