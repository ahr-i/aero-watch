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

	if err := s.createDriverInfoTable(); err != nil {
		return err
	}

	return s.createDroneDriverMatchTable()
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

func (s *MySQLStore) ListDroneModels() ([]DroneModel, error) {
	table, groupColumn, codeColumn, statusColumn, err := droneTableIdentifiers()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT %s, %s, %s FROM %s ORDER BY %s ASC, %s ASC",
		groupColumn,
		codeColumn,
		statusColumn,
		table,
		groupColumn,
		codeColumn,
	)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	drones := []DroneModel{}
	for rows.Next() {
		var drone DroneModel
		if err := rows.Scan(&drone.Group, &drone.Code, &drone.Status); err != nil {
			return nil, err
		}

		drones = append(drones, drone)
	}

	return drones, rows.Err()
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

func (s *MySQLStore) DeleteDroneModel(group string, code string) error {
	droneTable, groupColumn, codeColumn, _, err := droneTableIdentifiers()
	if err != nil {
		return err
	}

	matchTable, _, matchGroupColumn, matchCodeColumn, err := droneDriverMatchTableIdentifiers()
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	deleteMatchesQuery := fmt.Sprintf("DELETE FROM %s WHERE %s = ? AND %s = ?", matchTable, matchGroupColumn, matchCodeColumn)
	if _, err := tx.Exec(deleteMatchesQuery, group, code); err != nil {
		return err
	}

	deleteDroneQuery := fmt.Sprintf("DELETE FROM %s WHERE %s = ? AND %s = ?", droneTable, groupColumn, codeColumn)
	result, err := tx.Exec(deleteDroneQuery, group, code)
	if err != nil {
		return err
	}

	if err := noRowsErrorIfNotAffected(result); err != nil {
		return err
	}

	return tx.Commit()
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
	driverTable, idColumn, contentColumn, err := driverInfoTableIdentifiers()
	if err != nil {
		return nil, err
	}

	matchTable, matchDriverIDColumn, groupColumn, codeColumn, err := droneDriverMatchTableIdentifiers()
	if err != nil {
		return nil, err
	}

	droneTable, droneGroupColumn, droneCodeColumn, droneStatusColumn, err := droneTableIdentifiers()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT d.%s, d.%s, m.%s, m.%s, dr.%s
		FROM %s d
		LEFT JOIN %s m ON m.%s = d.%s
		LEFT JOIN %s dr ON dr.%s = m.%s AND dr.%s = m.%s
		ORDER BY d.%s ASC, m.%s ASC, m.%s ASC`,
		idColumn,
		contentColumn,
		groupColumn,
		codeColumn,
		droneStatusColumn,
		driverTable,
		matchTable,
		matchDriverIDColumn,
		idColumn,
		droneTable,
		droneGroupColumn,
		groupColumn,
		droneCodeColumn,
		codeColumn,
		idColumn,
		groupColumn,
		codeColumn,
	)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	infoByID := map[int64]*DriverInfo{}
	infoOrder := []int64{}
	for rows.Next() {
		var id int64
		var content string
		var group sql.NullString
		var code sql.NullString
		var status sql.NullString

		if err := rows.Scan(&id, &content, &group, &code, &status); err != nil {
			return nil, err
		}

		info, ok := infoByID[id]
		if !ok {
			info = &DriverInfo{
				ID:      id,
				Content: content,
				Drones:  []DroneModel{},
			}
			infoByID[id] = info
			infoOrder = append(infoOrder, id)
		}

		if group.Valid && code.Valid {
			info.Drones = append(info.Drones, DroneModel{
				Group:  group.String,
				Code:   code.String,
				Status: status.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	infos := make([]DriverInfo, 0, len(infoOrder))
	for _, id := range infoOrder {
		infos = append(infos, *infoByID[id])
	}

	return infos, nil
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

	matchTable, matchDriverIDColumn, _, _, err := droneDriverMatchTableIdentifiers()
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	deleteMatchesQuery := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", matchTable, matchDriverIDColumn)
	if _, err := tx.Exec(deleteMatchesQuery, id); err != nil {
		return err
	}

	deleteDriverQuery := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", table, idColumn)
	result, err := tx.Exec(deleteDriverQuery, id)
	if err != nil {
		return err
	}

	if err := noRowsErrorIfNotAffected(result); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *MySQLStore) CreateDroneDriverMatch(driverID int64, group string, code string) error {
	if err := s.ensureDriverExists(driverID); err != nil {
		return err
	}

	if err := s.ensureDroneExists(group, code); err != nil {
		return err
	}

	table, driverIDColumn, groupColumn, codeColumn, err := droneDriverMatchTableIdentifiers()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES (?, ?, ?)",
		table,
		driverIDColumn,
		groupColumn,
		codeColumn,
	)

	_, err = s.db.Exec(query, driverID, group, code)
	if isDuplicateMatchingError(err) {
		return ErrMatchingAlreadyExists
	}

	return err
}

func (s *MySQLStore) FindDriverInfoByDrone(group string, code string) (DriverInfo, error) {
	driverTable, driverIDColumn, contentColumn, err := driverInfoTableIdentifiers()
	if err != nil {
		return DriverInfo{}, err
	}

	matchTable, matchDriverIDColumn, groupColumn, codeColumn, err := droneDriverMatchTableIdentifiers()
	if err != nil {
		return DriverInfo{}, err
	}

	query := fmt.Sprintf(`SELECT d.%s, d.%s
		FROM %s m
		JOIN %s d ON d.%s = m.%s
		WHERE m.%s = ? AND m.%s = ?
		LIMIT 1`,
		driverIDColumn,
		contentColumn,
		matchTable,
		driverTable,
		driverIDColumn,
		matchDriverIDColumn,
		groupColumn,
		codeColumn,
	)

	var info DriverInfo
	err = s.db.QueryRow(query, group, code).Scan(&info.ID, &info.Content)
	if err != nil {
		return DriverInfo{}, err
	}

	return info, nil
}

func (s *MySQLStore) DeleteDroneDriverMatch(driverID int64, group string, code string) error {
	table, driverIDColumn, groupColumn, codeColumn, err := droneDriverMatchTableIdentifiers()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ? AND %s = ? AND %s = ?",
		table,
		driverIDColumn,
		groupColumn,
		codeColumn,
	)

	result, err := s.db.Exec(query, driverID, group, code)
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

func (s *MySQLStore) createDroneDriverMatchTable() error {
	table, driverIDColumn, groupColumn, codeColumn, err := droneDriverMatchTableIdentifiers()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		%s BIGINT NOT NULL,
		%s VARCHAR(255) NOT NULL,
		%s VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (%s, %s)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`, table, driverIDColumn, groupColumn, codeColumn, groupColumn, codeColumn)

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

func droneDriverMatchTableIdentifiers() (string, string, string, string, error) {
	table, err := sqlIdentifier(setting.Setting.DroneDriverMatchTable.Name)
	if err != nil {
		return "", "", "", "", err
	}

	driverIDColumn, err := sqlIdentifier(setting.Setting.DroneDriverMatchTable.DriverIDColumn)
	if err != nil {
		return "", "", "", "", err
	}

	groupColumn, err := sqlIdentifier(setting.Setting.DroneDriverMatchTable.GroupColumn)
	if err != nil {
		return "", "", "", "", err
	}

	codeColumn, err := sqlIdentifier(setting.Setting.DroneDriverMatchTable.CodeColumn)
	if err != nil {
		return "", "", "", "", err
	}

	return table, driverIDColumn, groupColumn, codeColumn, nil
}

func (s *MySQLStore) ensureDriverExists(id int64) error {
	table, idColumn, _, err := driverInfoTableIdentifiers()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT 1 FROM %s WHERE %s = ? LIMIT 1", table, idColumn)

	var exists int
	return s.db.QueryRow(query, id).Scan(&exists)
}

func (s *MySQLStore) ensureDroneExists(group string, code string) error {
	table, groupColumn, codeColumn, _, err := droneTableIdentifiers()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT 1 FROM %s WHERE %s = ? AND %s = ? LIMIT 1", table, groupColumn, codeColumn)

	var exists int
	return s.db.QueryRow(query, group, code).Scan(&exists)
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

func isDuplicateMatchingError(err error) bool {
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
