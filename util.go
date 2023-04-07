package zutils

import (
	"encoding/json"
	"github.com/google/uuid"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// FirstValue
// @Description: 返回第一个真值
// @param args
// @return T
func FirstValue[T Object](args ...T) T {
	for _, item := range args {
		// 跳过无效值
		if !reflect.ValueOf(item).IsValid() {
			continue
		}
		switch reflect.TypeOf(item).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
			reflect.Uintptr, reflect.Ptr:
			if !reflect.ValueOf(item).IsZero() {
				return item
			}
		case reflect.String:
			if reflect.ValueOf(item).String() != "" {
				return item
			}
		case reflect.Bool:
			if reflect.ValueOf(item).Bool() {
				return item
			}
		default:
			if !reflect.ValueOf(item).IsNil() {
				return item
			}
		}
	}
	// 如果无真值，返回第一个
	return args[0]
}

func StructToString(data interface{}) (string, error) {
	if str, err := json.Marshal(data); err != nil {
		return "", err
	} else {
		return string(str), nil
	}
}

// NewUuid
// @Auth: oak  2021-10-15 18:36:29
// @Description:  UUIDv4
// @return string
func NewUuid() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func DefaultPage(h interface{}) (int, int) {
	immutable := reflect.ValueOf(h)
	page := reflect.Indirect(immutable).FieldByName("Page").Int()
	rows := reflect.Indirect(immutable).FieldByName("Rows").Int()
	if page == 0 {
		page = 1
	}
	if rows == 0 {
		rows = 50
	}
	if rows > 1000 {
		rows = 1000
	}
	return int(page), int(rows)
}

// ScanPort
// @Description: 扫描端口
// @param protocol
// @param hostname
// @param port
// @return bool
func ScanPort(protocol string, hostname string, port int) bool {
	p := strconv.Itoa(port)
	addr := net.JoinHostPort(hostname, p)
	conn, err := net.DialTimeout(protocol, addr, 3*time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func Interface2Struct[T interface{}](src interface{}, dst *T) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}
