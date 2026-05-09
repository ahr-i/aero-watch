package db

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ahr-i/aero-watch/emergency/setting"
	_ "github.com/go-sql-driver/mysql"
)

var identifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type MySQLController struct {
	db *sql.DB
}

func NewMySQLController() (Controller, error) {
	cfg := setting.Setting.DB
	host := os.Getenv(cfg.HostEnv)
	port := os.Getenv(cfg.PortEnv)
	user := os.Getenv(cfg.UserEnv)
	password := os.Getenv(cfg.PasswordEnv)
	schema := os.Getenv(cfg.SchemaEnv)

	if host == "" || port == "" || user == "" || schema == "" {
		return nil, errors.New("database env values are missing")
	}

	host, port = normalizeHostPort(host, port)
	err := ensureSchema(host, port, user, password, schema)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", user, password, net.JoinHostPort(host, port), schema)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, err
	}

	return &MySQLController{db: sqlDB}, nil
}

func ensureSchema(host string, port string, user string, password string, schema string) error {
	if err := validateSchemaName(schema); err != nil {
		return err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=true&loc=Local", user, password, net.JoinHostPort(host, port))
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", quoteIdentifier(schema)))
	return err
}

func (c *MySQLController) Close() error {
	return c.db.Close()
}

func (c *MySQLController) ImportCSV(tableName string, columns []string, rows [][]string) error {
	if err := validateTableName(tableName); err != nil {
		return err
	}
	if err := validateColumns(columns); err != nil {
		return err
	}

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	definitions := make([]string, len(columns))
	for i, column := range columns {
		definitions[i] = fmt.Sprintf("%s TEXT", quoteIdentifier(column))
	}

	_, err = tx.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", quoteIdentifier(tableName), strings.Join(definitions, ", ")))
	if err != nil {
		return err
	}

	_, err = tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", quoteIdentifier(tableName)))
	if err != nil {
		return err
	}

	if len(rows) > 0 {
		quotedColumns := make([]string, len(columns))
		placeholders := make([]string, len(columns))
		for i, column := range columns {
			quotedColumns[i] = quoteIdentifier(column)
			placeholders[i] = "?"
		}

		stmt, err := tx.Prepare(fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s)",
			quoteIdentifier(tableName),
			strings.Join(quotedColumns, ", "),
			strings.Join(placeholders, ", "),
		))
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, row := range rows {
			values := make([]any, len(row))
			for i, value := range row {
				values[i] = value
			}
			if _, err := stmt.Exec(values...); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (c *MySQLController) DropTable(tableName string) error {
	if err := validateTableName(tableName); err != nil {
		return err
	}

	_, err := c.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", quoteIdentifier(tableName)))
	return err
}

func (c *MySQLController) ListTables() ([]string, error) {
	prefix := setting.Setting.TablePrefix + "_"
	rows, err := c.db.Query(
		"SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name LIKE ? ORDER BY table_name",
		escapeLike(prefix)+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tablePattern := regexp.MustCompile("^" + regexp.QuoteMeta(prefix) + `[A-Za-z][A-Za-z0-9]*_([0-9]{6}|[0-9]{8})$`)
	tables := make([]string, 0)
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		if tablePattern.MatchString(tableName) {
			tables = append(tables, tableName)
		}
	}

	return tables, rows.Err()
}

func (c *MySQLController) GetTable(tableName string) (*TableData, error) {
	if err := validateTableName(tableName); err != nil {
		return nil, err
	}

	rows, err := c.db.Query(fmt.Sprintf("SELECT * FROM %s", quoteIdentifier(tableName)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rawColumns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columns := make([]string, len(rawColumns))
	for i, column := range rawColumns {
		columns[i] = snakeToCamel(column)
	}

	data := make([]map[string]string, 0)
	for rows.Next() {
		values := make([]sql.NullString, len(rawColumns))
		scanValues := make([]any, len(rawColumns))
		for i := range values {
			scanValues[i] = &values[i]
		}

		if err := rows.Scan(scanValues...); err != nil {
			return nil, err
		}

		row := make(map[string]string, len(rawColumns))
		for i, column := range rawColumns {
			if values[i].Valid {
				row[snakeToCamel(column)] = values[i].String
			} else {
				row[snakeToCamel(column)] = ""
			}
		}
		data = append(data, row)
	}

	return &TableData{
		Columns: columns,
		Rows:    data,
	}, rows.Err()
}

func (c *MySQLController) FindNearest(latitude float64, longitude float64, limit int) ([]NearbyEmergency, error) {
	if limit < 1 {
		limit = 1
	}

	tables, err := c.ListTables()
	if err != nil {
		return nil, err
	}

	nearby := make([]NearbyEmergency, 0)
	for _, tableName := range tables {
		tableData, err := c.GetTable(tableName)
		if err != nil {
			return nil, err
		}

		for _, row := range tableData.Rows {
			rowLongitude, err := strconv.ParseFloat(row["longitude"], 64)
			if err != nil {
				continue
			}

			rowLatitude, err := strconv.ParseFloat(row["latitude"], 64)
			if err != nil {
				continue
			}

			nearby = append(nearby, NearbyEmergency{
				TableName:     tableName,
				DistanceMeter: distanceMeter(latitude, longitude, rowLatitude, rowLongitude),
				Data:          row,
			})
		}
	}

	sort.Slice(nearby, func(i int, j int) bool {
		return nearby[i].DistanceMeter < nearby[j].DistanceMeter
	})

	if len(nearby) > limit {
		nearby = nearby[:limit]
	}

	return nearby, nil
}

func distanceMeter(latitude1 float64, longitude1 float64, latitude2 float64, longitude2 float64) float64 {
	const earthRadiusMeter = 6371000

	lat1 := latitude1 * math.Pi / 180
	lat2 := latitude2 * math.Pi / 180
	deltaLat := (latitude2 - latitude1) * math.Pi / 180
	deltaLon := (longitude2 - longitude1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeter * c
}

func snakeToCamel(value string) string {
	parts := strings.Split(value, "_")
	if len(parts) == 1 {
		return value
	}

	builder := strings.Builder{}
	builder.WriteString(parts[0])
	for _, part := range parts[1:] {
		if part == "" {
			continue
		}
		builder.WriteString(strings.ToUpper(part[:1]))
		if len(part) > 1 {
			builder.WriteString(part[1:])
		}
	}

	return builder.String()
}

func validateTableName(tableName string) error {
	if !identifierPattern.MatchString(tableName) {
		return fmt.Errorf("invalid table name: %s", tableName)
	}
	return nil
}

func validateSchemaName(schema string) error {
	if !identifierPattern.MatchString(schema) {
		return fmt.Errorf("invalid schema name: %s", schema)
	}
	return nil
}

func validateColumns(columns []string) error {
	if len(columns) == 0 {
		return errors.New("columns are required")
	}

	for _, column := range columns {
		if !identifierPattern.MatchString(column) {
			return fmt.Errorf("invalid column name: %s", column)
		}
	}

	return nil
}

func quoteIdentifier(identifier string) string {
	return "`" + identifier + "`"
}

func escapeLike(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `%`, `\%`)
	return strings.ReplaceAll(value, `_`, `\_`)
}

func normalizeHostPort(host string, port string) (string, string) {
	parsed, err := url.Parse(host)
	if err == nil && parsed.Host != "" {
		host = parsed.Host
	}

	if h, p, err := net.SplitHostPort(host); err == nil {
		return h, p
	}

	if strings.Contains(host, ":") && !strings.Contains(host, "]") {
		parts := strings.Split(host, ":")
		if len(parts) == 2 {
			return parts[0], parts[1]
		}
	}

	return host, port
}
