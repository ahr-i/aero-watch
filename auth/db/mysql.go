package db

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/ahr-i/aero-watch/auth/setting"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore() Store {
	return &MySQLStore{}
}

func (s *MySQLStore) Init() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		setting.Setting.DB.User,
		setting.Setting.DB.Password,
		setting.Setting.DB.Host,
		setting.Setting.DB.Port,
		setting.Setting.DB.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return err
	}

	s.db = db
	if err := s.createUserTable(); err != nil {
		return err
	}

	return s.createAdminUserIfNotExists()
}

func (s *MySQLStore) CreateUser(user string, passwordHash string, role string) error {
	table, err := sqlIdentifier(setting.Setting.UserTable.Name)
	if err != nil {
		return err
	}

	usernameColumn, err := sqlIdentifier(setting.Setting.UserTable.UsernameColumn)
	if err != nil {
		return err
	}

	passwordHashColumn, err := sqlIdentifier(setting.Setting.UserTable.PasswordHashColumn)
	if err != nil {
		return err
	}

	roleColumn, err := sqlIdentifier(setting.Setting.UserTable.RoleColumn)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES (?, ?, ?)",
		table,
		usernameColumn,
		passwordHashColumn,
		roleColumn,
	)

	_, err = s.db.Exec(query, user, passwordHash, role)
	if isDuplicateUserError(err) {
		return ErrUserAlreadyExists
	}

	return err
}

func (s *MySQLStore) FindUserAuthInfo(user string) (UserAuthInfo, error) {
	table, err := sqlIdentifier(setting.Setting.UserTable.Name)
	if err != nil {
		return UserAuthInfo{}, err
	}

	usernameColumn, err := sqlIdentifier(setting.Setting.UserTable.UsernameColumn)
	if err != nil {
		return UserAuthInfo{}, err
	}

	passwordHashColumn, err := sqlIdentifier(setting.Setting.UserTable.PasswordHashColumn)
	if err != nil {
		return UserAuthInfo{}, err
	}

	roleColumn, err := sqlIdentifier(setting.Setting.UserTable.RoleColumn)
	if err != nil {
		return UserAuthInfo{}, err
	}

	query := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s = ? LIMIT 1",
		passwordHashColumn,
		roleColumn,
		table,
		usernameColumn,
	)

	var userInfo UserAuthInfo
	err = s.db.QueryRow(query, user).Scan(&userInfo.PasswordHash, &userInfo.Role)
	if err != nil {
		return UserAuthInfo{}, err
	}

	return userInfo, nil
}

func (s *MySQLStore) ListUsers() ([]UserInfo, error) {
	table, err := sqlIdentifier(setting.Setting.UserTable.Name)
	if err != nil {
		return nil, err
	}

	usernameColumn, err := sqlIdentifier(setting.Setting.UserTable.UsernameColumn)
	if err != nil {
		return nil, err
	}

	roleColumn, err := sqlIdentifier(setting.Setting.UserTable.RoleColumn)
	if err != nil {
		return nil, err
	}

	createdAtColumn, err := sqlIdentifier(setting.Setting.UserTable.CreatedAtColumn)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT %s, %s, %s FROM %s ORDER BY %s DESC",
		usernameColumn,
		roleColumn,
		createdAtColumn,
		table,
		createdAtColumn,
	)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []UserInfo{}
	for rows.Next() {
		var user UserInfo
		if err := rows.Scan(&user.User, &user.Role, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *MySQLStore) UpdateUserRole(user string, role string) error {
	table, err := sqlIdentifier(setting.Setting.UserTable.Name)
	if err != nil {
		return err
	}

	usernameColumn, err := sqlIdentifier(setting.Setting.UserTable.UsernameColumn)
	if err != nil {
		return err
	}

	roleColumn, err := sqlIdentifier(setting.Setting.UserTable.RoleColumn)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s = ?", table, roleColumn, usernameColumn)
	result, err := s.db.Exec(query, role, user)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *MySQLStore) DeleteUser(user string) error {
	table, err := sqlIdentifier(setting.Setting.UserTable.Name)
	if err != nil {
		return err
	}

	usernameColumn, err := sqlIdentifier(setting.Setting.UserTable.UsernameColumn)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", table, usernameColumn)
	result, err := s.db.Exec(query, user)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *MySQLStore) Close() error {
	if s.db == nil {
		return nil
	}

	return s.db.Close()
}

func (s *MySQLStore) createUserTable() error {
	table, err := sqlIdentifier(setting.Setting.UserTable.Name)
	if err != nil {
		return err
	}

	usernameColumn, err := sqlIdentifier(setting.Setting.UserTable.UsernameColumn)
	if err != nil {
		return err
	}

	passwordHashColumn, err := sqlIdentifier(setting.Setting.UserTable.PasswordHashColumn)
	if err != nil {
		return err
	}

	roleColumn, err := sqlIdentifier(setting.Setting.UserTable.RoleColumn)
	if err != nil {
		return err
	}

	createdAtColumn, err := sqlIdentifier(setting.Setting.UserTable.CreatedAtColumn)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		%s VARCHAR(255) NOT NULL PRIMARY KEY,
		%s VARCHAR(255) NOT NULL,
		%s VARCHAR(50) NOT NULL,
		%s TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`, table, usernameColumn, passwordHashColumn, roleColumn, createdAtColumn)

	_, err = s.db.Exec(query)
	return err
}

func (s *MySQLStore) createAdminUserIfNotExists() error {
	userInfo, err := s.FindUserAuthInfo(setting.Setting.Admin.User)
	if err == nil {
		if userInfo.Role == setting.Setting.Role.Admin {
			return nil
		}

		return s.UpdateUserRole(setting.Setting.Admin.User, setting.Setting.Role.Admin)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(setting.Setting.Admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.CreateUser(setting.Setting.Admin.User, string(passwordHash), setting.Setting.Role.Admin)
}

func isDuplicateUserError(err error) bool {
	if err == nil {
		return false
	}

	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

var sqlIdentifierRegexp = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

func sqlIdentifier(value string) (string, error) {
	if !sqlIdentifierRegexp.MatchString(value) {
		return "", fmt.Errorf("invalid SQL identifier: %s", value)
	}

	return "`" + value + "`", nil
}
