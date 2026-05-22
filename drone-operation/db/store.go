package db

import (
	"errors"
	"fmt"

	"github.com/ahr-i/aero-watch/drone-operation/setting"
)

var ErrDroneAlreadyExists = errors.New("drone model already exists")
var ErrInvalidStatus = errors.New("invalid drone status")
var ErrMatchingAlreadyExists = errors.New("matching already exists")

type Store interface {
	Init() error
	RegisterDroneModel(group string, code string) error
	ListDroneModels() ([]DroneModel, error)
	ValidateDroneModel(group string, code string) error
	UpdateDroneStatus(group string, code string, status string) error
	DeleteDroneModel(group string, code string) error
	CreateDriverInfo(content string) (DriverInfo, error)
	ListDriverInfos() ([]DriverInfo, error)
	UpdateDriverInfo(id int64, content string) error
	DeleteDriverInfo(id int64) error
	CreateDroneDriverMatch(driverID int64, group string, code string) error
	FindDriverInfoByDrone(group string, code string) (DriverInfo, error)
	DeleteDroneDriverMatch(driverID int64, group string, code string) error
	Close() error
}

type DriverInfo struct {
	ID      int64
	Content string
	Drones  []DroneModel
}

type DroneModel struct {
	Group  string
	Code   string
	Status string
}

func NewStore() (Store, error) {
	switch setting.Setting.DB.Type {
	case "mysql":
		return NewMySQLStore(), nil
	default:
		return nil, fmt.Errorf("unsupported db type: %s", setting.Setting.DB.Type)
	}
}
