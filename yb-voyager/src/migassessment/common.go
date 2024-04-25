/*
Copyright (c) YugabyteDB, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package migassessment

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var AssessmentMetadataDir string
var TargetYBVersion string

type Record map[string]any

var SizingReport *AssessmentReport

type AssessmentReport struct {
	ColocatedTables                 []string
	ColocatedReasoning              string
	ShardedTables                   []string
	NumNodes                        float64
	VCPUsPerInstance                float64
	MemoryPerInstance               float64
	OptimalSelectConnectionsPerNode int64
	OptimalInsertConnectionsPerNode int64
	MigrationTimeTakenInMin         float64
	// ParallelVoyagerImportThreads int64 ==> Optional
}

var assessmentParams = &AssessmentParams{}

type AssessmentParams struct {
	TargetYBVersion string `toml:"target_yb_version"`
}

func LoadCSVDataFile[T any](filePath string) ([]*T, error) {
	result := make([]*T, 0)
	records, err := loadCSVDataFileGeneric(filePath)
	if err != nil {
		log.Errorf("error loading csv data file %s: %v", filePath, err)
		return nil, fmt.Errorf("error loading csv data file %s: %w", filePath, err)
	}

	for _, record := range records {
		var tmplRec T
		bs, err := json.Marshal(record)
		if err != nil {
			log.Errorf("error marshalling record: %v", err)
			return nil, fmt.Errorf("error marshalling record: %w", err)
		}
		err = json.Unmarshal(bs, &tmplRec)
		if err != nil {
			log.Errorf("error unmarshalling record: %v", err)
			return nil, fmt.Errorf("error unmarshalling record: %w", err)
		}
		result = append(result, &tmplRec)
	}
	return result, nil
}

func loadCSVDataFileGeneric(filePath string) ([]Record, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Warnf("error opening file %s: %v", filePath, err)
		return nil, nil
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Errorf("error closing file %s: %v", filePath, err)
		}
	}()

	csvReader := csv.NewReader(file)
	csvReader.ReuseRecord = true

	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Errorf("error reading csv file %s: %v", filePath, err)
		return nil, fmt.Errorf("error reading csv file %s: %w", filePath, err)
	}

	if len(rows) == 0 {
		log.Warnf("file '%s' is empty, no records", filePath)
		return nil, nil
	}

	columnNames := rows[0]
	result := make([]Record, len(rows)-1)
	for rowNum := 1; rowNum < len(rows); rowNum++ {
		record := make(Record)
		row := rows[rowNum]
		for columnIdx, columnName := range columnNames {
			record[columnName] = row[columnIdx]
		}
		result[rowNum-1] = record
	}
	return result, nil
}

func LoadAssessmentParams(userInputFpath string) error {
	if userInputFpath == "" {
		log.Infof("user input file path is empty, skipping loading assessment parameters")
		return nil
	}
	log.Infof("loading assessment parameters from file %s", userInputFpath)
	tomlData, err := os.ReadFile(userInputFpath)
	if err != nil {
		log.Errorf("error reading toml file %s: %v", userInputFpath, err)
		return fmt.Errorf("error reading toml file %s: %w", userInputFpath, err)
	}

	err = toml.Unmarshal(tomlData, &assessmentParams)
	if err != nil {
		log.Errorf("error unmarshalling toml file's data: %v", err)
		return fmt.Errorf("error unmarshalling toml file's data: %w", err)
	}

	return nil
}

func convertToMap(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	var allMaps []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		err := rows.Scan(pointers...)
		if err != nil {
			panic(err)
		}
		resultMap := make(map[string]interface{})
		for i, val := range values {
			//fmt.Printf("Adding key=%s val=%v\n", columns[i], val)
			resultMap[columns[i]] = val
		}
		allMaps = append(allMaps, resultMap)
	}
	return allMaps
}

func checkInternetAccess() (ok bool) {
	_, err := http.Get("http://clients3.google.com/generate_204")
	return err == nil
}
