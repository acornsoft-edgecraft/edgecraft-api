package utils

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// GetValuesFromInterface - 지정된 Interface 형식의 Structure에서 지정한 FIeld 이름을 기준으로 값을 Array로 반환 (using Reflect)
func GetValuesFromInterface(val interface{}, fields ...string) []interface{} {
	returnArray := make([]interface{}, 0)
	structVal := reflect.ValueOf(val).Elem()
	for _, name := range fields {
		field := structVal.FieldByName(name).Interface()
		returnArray = append(returnArray, field)
	}
	return returnArray
}

// GetDataFromRows - Structure Mapping 없이 SQL Result 처리
func GetDataFromRows(rows *sql.Rows) ([]interface{}, error) {
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	count := len(colTypes)
	finalRows := []interface{}{}

	for rows.Next() {
		scanArgs := make([]interface{}, count)
		for i, v := range colTypes {
			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP", "TIMESTAMPTZ":
				scanArgs[i] = new(sql.NullString)
				break
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break
			case "INT4":
				scanArgs[i] = new(sql.NullInt64)
				break
			case "NUMERIC":
				scanArgs[i] = new(sql.NullFloat64)
			default:
				scanArgs[i] = new(sql.NullString)
				break
			}
		}

		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		masterData := map[string]interface{}{}

		for i, v := range colTypes {
			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				masterData[v.Name()] = z.Bool
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				masterData[v.Name()] = z.Int64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				masterData[v.Name()] = z.Float64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				masterData[v.Name()] = z.Int32
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				masterData[v.Name()] = z.String
				continue
			}

			masterData[v.Name()] = scanArgs[i]
		}

		finalRows = append(finalRows, masterData)
	}

	return finalRows, nil
}

func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}

//RenderTmpl asdf
func RenderTmpl(tmplName string, templ string, obj interface{}) (string, error) {
	var commonTemplate *template.Template
	var muLock sync.RWMutex // 생성시 중복 방지를 위한 Lock

	if commonTemplate == nil {
		muLock.Lock()
		if commonTemplate == nil {
			funcMap := template.FuncMap{
				"isNull": func(val interface{}) bool {
					if val == nil || (reflect.ValueOf(val).Kind() == reflect.Ptr && reflect.ValueOf(val).IsNil()) {
						return false
					}
					return true
				},
			}
			commonTemplate = template.New(tmplName).Funcs(funcMap)
			fmt.Println("New template : ", commonTemplate)
		}
		muLock.Unlock()
	}

	t, err := commonTemplate.Parse(templ)

	var buff *bytes.Buffer
	if err == nil {
		buff = bytes.NewBufferString("")
		err = t.Execute(buff, obj)
	}

	return buff.String(), err
}

func ChkPortOpen(host string, port int) (bool, error) {

	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), timeout)
	if err != nil {
		return false, err
	}
	if conn != nil {
		defer conn.Close()
		return true, err
	}
	return false, err
}

func Print(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b), "\n")
}
