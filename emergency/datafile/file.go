package datafile

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/ahr-i/aero-watch/emergency/setting"
)

type DataFile struct {
	Name     string `json:"name"`
	Date     string `json:"date"`
	FileName string `json:"fileName"`
}

type CSVData struct {
	Columns []string   `json:"columns"`
	Rows    [][]string `json:"rows"`
}

func ValidateFiles() error {
	files, err := ListFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		_, err := ReadCSV(file.Name, file.Date)
		if err != nil {
			return err
		}
	}

	return nil
}

func ListFiles() ([]DataFile, error) {
	entries, err := os.ReadDir(setting.Setting.CSVPath)
	if err != nil {
		return nil, err
	}

	pattern, err := regexp.Compile(setting.Setting.DataFilePattern)
	if err != nil {
		return nil, err
	}

	files := make([]DataFile, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			return nil, fmt.Errorf("invalid csv entry: %s", entry.Name())
		}

		name := entry.Name()
		if _, ok := setting.Setting.CSVSchema[name]; !ok {
			return nil, fmt.Errorf("schema is not configured for csv directory: %s", name)
		}

		dateFiles, err := os.ReadDir(filepath.Join(setting.Setting.CSVPath, name))
		if err != nil {
			return nil, err
		}

		for _, dateFile := range dateFiles {
			if dateFile.IsDir() {
				return nil, fmt.Errorf("invalid csv date entry: %s/%s", name, dateFile.Name())
			}

			fileName := dateFile.Name()
			if filepath.Ext(fileName) != ".csv" {
				return nil, fmt.Errorf("invalid csv file extension: %s/%s", name, fileName)
			}

			pathName := filepath.ToSlash(filepath.Join(name, fileName))
			matches := pattern.FindStringSubmatch(pathName)
			if len(matches) != 3 {
				return nil, fmt.Errorf("invalid csv file name: %s", pathName)
			}

			files = append(files, DataFile{
				Name:     matches[1],
				Date:     matches[2],
				FileName: pathName,
			})
		}
	}

	return files, nil
}

func ReadCSV(name string, date string) (*CSVData, error) {
	if err := ValidateNameDate(name, date); err != nil {
		return nil, err
	}

	path := filepath.Join(setting.Setting.CSVPath, name, fmt.Sprintf("%s.csv", date))
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, errors.New("csv file is empty")
	}

	columns := records[0]
	expectedColumns, ok := setting.Setting.CSVSchema[name]
	if !ok {
		return nil, fmt.Errorf("schema is not configured for name: %s", name)
	}
	if !slices.Equal(columns, expectedColumns) {
		return nil, fmt.Errorf("csv columns do not match configured schema for name: %s", name)
	}

	for i, row := range records[1:] {
		if len(row) != len(columns) {
			return nil, fmt.Errorf("csv row %d has %d columns, want %d", i+2, len(row), len(columns))
		}
	}

	return &CSVData{
		Columns: columns,
		Rows:    records[1:],
	}, nil
}

func ValidateNameDate(name string, date string) error {
	if name == "" || date == "" {
		return errors.New("name and date are required")
	}

	pattern, err := regexp.Compile(setting.Setting.DataFilePattern)
	if err != nil {
		return err
	}
	if !pattern.MatchString(fmt.Sprintf("%s/%s.csv", name, date)) {
		return fmt.Errorf("invalid name or date: %s/%s", name, date)
	}

	if strings.Contains(name, "_") {
		return errors.New("name cannot contain underscore")
	}

	return nil
}

func TableName(name string, date string) string {
	return fmt.Sprintf("%s_%s_%s", setting.Setting.TablePrefix, name, date)
}
