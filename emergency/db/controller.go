package db

type Controller interface {
	Close() error
	ImportCSV(tableName string, columns []string, rows [][]string) error
	DropTable(tableName string) error
	ListTables() ([]string, error)
	GetTable(tableName string) (*TableData, error)
	FindNearest(latitude float64, longitude float64, limit int) ([]NearbyEmergency, error)
}

type TableData struct {
	Columns []string            `json:"columns"`
	Rows    []map[string]string `json:"rows"`
}

type NearbyEmergency struct {
	TableName     string            `json:"tableName"`
	DistanceMeter float64           `json:"distanceMeter"`
	Data          map[string]string `json:"data"`
}
