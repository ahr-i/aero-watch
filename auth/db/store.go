package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/ahr-i/aero-watch/auth/setting"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type Store interface {
	Init() error
	CreateUser(user string, passwordHash string, role string) error
	FindUserAuthInfo(user string) (UserAuthInfo, error)
	ListUsers() ([]UserInfo, error)
	UpdateUserRole(user string, role string) error
	DeleteUser(user string) error
	Close() error
}

type UserAuthInfo struct {
	PasswordHash string
	Role         string
}

type UserInfo struct {
	User      string
	Role      string
	CreatedAt time.Time
}

func NewStore() (Store, error) {
	switch setting.Setting.DB.Type {
	case "mysql":
		return NewMySQLStore(), nil
	default:
		return nil, fmt.Errorf("unsupported db type: %s", setting.Setting.DB.Type)
	}
}
