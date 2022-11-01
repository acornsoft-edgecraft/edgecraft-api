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
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
)

// getTypeString - 지정한 Go 형식을 Refelct를 이용해서 형식 문자열로 반환
func getTypeString(t reflect.Type) string {
	if t.PkgPath() == "main" {
		return t.Name()
	}
	return t.String()
}

// getGoString - 지정한 형식의 데이터를 Reflect를 이용해서 문자열로 반환
func getGoString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "nil"
	case reflect.Struct:
		t := v.Type()
		out := getTypeString(t) + "{"
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				out += ", "
			}
			fieldValue := v.Field(i)
			field := t.Field(i)
			out += fmt.Sprintf("%s: %s", field.Name, getGoString(fieldValue))
		}
		out += "}"
		return out
	case reflect.Interface, reflect.Ptr:
		if v.IsZero() {
			return fmt.Sprintf("(%s)(nil)", getTypeString(v.Type()))
		}
		return "&" + getGoString(v.Elem())
	case reflect.Slice:
		out := getTypeString(v.Type())
		if v.IsZero() {
			out += "(nil)"
		} else {
			out += "{"
			for i := 0; i < v.Len(); i++ {
				if i > 0 {
					out += ", "
				}
				out += getGoString(v.Index(i))
			}
			out += "}"
		}
		return out
	default:
		return fmt.Sprintf("%#v", v)
	}
}

// GetGoString - Go Struct를 문자열로 반환
func GetGoString(v interface{}) string {
	return getGoString(reflect.ValueOf(v))
}

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
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
			case "INT4":
				scanArgs[i] = new(sql.NullInt64)
			case "NUMERIC":
				scanArgs[i] = new(sql.NullFloat64)
			default:
				scanArgs[i] = new(sql.NullString)
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

// RenderTmpl asdf
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

// GenerateUUID - UUID 생성
func GenerateUUID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

// StringPtr - 지정한 문자열을 포인터로 반환
func StringPtr(val string) *string {
	return &val
}

// TimePtr - 지정한 Time정보를 포인터로 반환
func TimePtr(val time.Time) *time.Time {
	return &val
}

// IntPrt - 지정한 int 정보를 포인터로 반환
func IntPrt(val int) *int {
	return &val
}

// BoolPtr - 지정한 Bool 정보를 포인터로 반환
func BoolPtr(val bool) *bool {
	return &val
}

// ArrayContains - 지정한 배열내에 지정한 값이 있는지 검증
func ArrayContains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

// EndWithOnArray - 지정한 배열내의 값들 중에서 EndsWith에 해당하는 값이 있는지 검증
func EndWithOnArray(arr []string, str string) string {
	for _, s := range arr {
		if strings.HasSuffix(s, str) {
			return s
		}
	}
	return ""
}
