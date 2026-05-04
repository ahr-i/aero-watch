package db

import (
	"errors"
	"fmt"

	"github.com/ahr-i/aero-watch/auth/setting"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type Store interface {
	Init() error
	CreateUser(user string, passwordHash string, role string) error
	FindUserAuthInfo(user string) (UserAuthInfo, error)
	UpdateUserRole(user string, role string) error
	Close() error
}

type UserAuthInfo struct {
	PasswordHash string
	Role         string
}

func NewStore() (Store, error) {
	switch setting.Setting.DB.Type {
	case "mysql":
		return NewMySQLStore(), nil
	default:
		return nil, fmt.Errorf("unsupported db type: %s", setting.Setting.DB.Type)
	}
}
