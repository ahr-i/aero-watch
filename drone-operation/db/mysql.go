package db

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/ahr-i/aero-watch/drone-operation/setting"
	"github.com/go-sql-driver/mysql"
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
	if err := s.createDroneTable(); err != nil {
		return err
	}

	return s.createDriverInfoTable()
}

func (s *MySQLStore) RegisterDroneModel(group string, code string) error {
	table, groupColumn, codeColumn, statusColumn, err := droneTableIdentifiers()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES (?, ?, ?)",
		table,
		groupColumn,
		codeColumn,
		statusColumn,
	)

	_, err = s.db.Exec(query, group, code, setting.Setting.Status.Active)
	if isDuplicateDroneError(err) {
		return ErrDroneAlreadyExists
	}

	return err
}

func (s *MySQLStore) ValidateDroneModel(group string, code string) error {
	table, groupColumn, codeColumn, statusColumn, err := droneTableIdentifiers()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT 1 FROM %s WHERE %s = ? AND %s = ? AND %s = ? LIMIT 1",
		table,
		groupColumn,
		codeColumn,
		statusColumn,
	)

	var exists int
	err = s.db.QueryRow(query, group, code, setting.Setting.Status.Active).Scan(&exists)
	return err
}

func (s *MySQLStore) UpdateDroneStatus(group string, code string, status string) error {
	if status != setting.Setting.Status.Active && status != setting.Setting.Status.Inactive {
		return ErrInvalidStatus
	}

	table, groupColumn, codeColumn, statusColumn, err := droneTableIdentifiers()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s = ? AND %s = ?",
		table,
		statusColumn,
		groupColumn,
		codeColumn,
	)

	result, err := s.db.Exec(query, status, group, code)
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

func (s *MySQLStore) CreateDriverInfo(content string) (DriverInfo, error) {
	table, _, contentColumn, err := driverInfoTableIdentifiers()
	if err != nil {
		return DriverInfo{}, err
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (?)", table, contentColumn)
	result, err := s.db.Exec(query, content)
	if err != nil {
		return DriverInfo{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return DriverInfo{}, err
	}

	return DriverInfo{
		ID:      id,
		Content: content,
	}, nil
}

func (s *MySQLStore) ListDriverInfos() ([]DriverInfo, error) {
	table, idColumn, contentColumn, err := driverInfoTableIdentifiers()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT %s, %s FROM %s ORDER BY %s ASC", idColumn, contentColumn, table, idColumn)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	infos := []DriverInfo{}
	for rows.Next() {
		var info DriverInfo
		if err := rows.Scan(&info.ID, &info.Content); err != nil {
			return nil, err
		}

		infos = append(infos, info)
	}

	return infos, rows.Err()
}

func (s *MySQLStore) UpdateDriverInfo(id int64, content string) error {
	table, idColumn, contentColumn, err := driverInfoTableIdentifiers()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s = ?", table, contentColumn, idColumn)
	result, err := s.db.Exec(query, content, id)
	if err != nil {
		return err
	}

	return noRowsErrorIfNotAffected(result)
}

func (s *MySQLStore) DeleteDriverInfo(id int64) error {
	table, idColumn, _, err := driverInfoTableIdentifiers()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", table, idColumn)
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return noRowsErrorIfNotAffected(result)
}

func (s *MySQLStore) Close() error {
	if s.db == nil {
		return nil
	}

	return s.db.Close()
}

func (s *MySQLStore) createDroneTable() error {
	table, groupColumn, codeColumn, statusColumn, err := droneTableIdentifiers()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		%s VARCHAR(255) NOT NULL,
		%s VARCHAR(255) NOT NULL,
		%s VARCHAR(50) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (%s, %s)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`, table, groupColumn, codeColumn, statusColumn, groupColumn, codeColumn)

	_, err = s.db.Exec(query)
	return err
}

func (s *MySQLStore) createDriverInfoTable() error {
	table, idColumn, contentColumn, err := driverInfoTableIdentifiers()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		%s BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		%s TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`, table, idColumn, contentColumn)

	_, err = s.db.Exec(query)
	return err
}

func droneTableIdentifiers() (string, string, string, string, error) {
	table, err := sqlIdentifier(setting.Setting.DroneTable.Name)
	if err != nil {
		return "", "", "", "", err
	}

	groupColumn, err := sqlIdentifier(setting.Setting.DroneTable.GroupColumn)
	if err != nil {
		return "", "", "", "", err
	}

	codeColumn, err := sqlIdentifier(setting.Setting.DroneTable.CodeColumn)
	if err != nil {
		return "", "", "", "", err
	}

	statusColumn, err := sqlIdentifier(setting.Setting.DroneTable.StatusColumn)
	if err != nil {
		return "", "", "", "", err
	}

	return table, groupColumn, codeColumn, statusColumn, nil
}

func driverInfoTableIdentifiers() (string, string, string, error) {
	table, err := sqlIdentifier(setting.Setting.DriverInfoTable.Name)
	if err != nil {
		return "", "", "", err
	}

	idColumn, err := sqlIdentifier(setting.Setting.DriverInfoTable.IDColumn)
	if err != nil {
		return "", "", "", err
	}

	contentColumn, err := sqlIdentifier(setting.Setting.DriverInfoTable.ContentColumn)
	if err != nil {
		return "", "", "", err
	}

	return table, idColumn, contentColumn, nil
}

func noRowsErrorIfNotAffected(result sql.Result) error {
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func isDuplicateDroneError(err error) bool {
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
