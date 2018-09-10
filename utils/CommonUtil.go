package utils

import (
	"database/sql"
	"encoding/json"
)

func Rowtojson(rows *sql.Rows) []map[string]interface{} {
	columns, err := rows.Columns()
	if err != nil {
		//fmt.Println("row to json:", err)
		return nil
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		//fmt.Println("2 row to json:", err)
		return nil
	}

	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(string(jsonData)), &result); err != nil {
		println(err.Error())
	}
	//fmt.Println(string(jsonData))
	return result
}
