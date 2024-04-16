package file

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// DataRow holds a single row of the file.
type DataRow struct {
	Name     string
	Revision string
}

func ReadCSV(filename string) ([]DataRow, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var data []DataRow

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading the file: %w", err)
		}

		// Trim leading and trailing spaces from each field
		trimmedName := strings.TrimSpace(record[0])
		trimmedRevision := strings.TrimSpace(record[1])

		data = append(data, DataRow{
			Name:     trimmedName,
			Revision: trimmedRevision,
		})
	}

	return data, nil
}
