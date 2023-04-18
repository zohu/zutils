package zutils

import (
	"encoding/json"
	"net"
	"reflect"
	"strconv"
	"time"
)

func FirstTruthString(args ...string) string {
	for _, item := range args {
		if item != "" {
			return item
		}
	}
	return args[0]
}
func FirstTruthInt(args ...int) int {
	for _, item := range args {
		if item != 0 {
			return item
		}
	}
	return args[0]
}

func FirstTruthInt64(args ...int64) int64 {
	for _, item := range args {
		if item != 0 {
			return item
		}
	}
	return args[0]
}

func FirstTruthFloat64(args ...float64) float64 {
	for _, item := range args {
		if item != 0 {
			return item
		}
	}
	return args[0]
}

func StructToString(data interface{}) (string, error) {
	if str, err := json.Marshal(data); err != nil {
		return "", err
	} else {
		return string(str), nil
	}
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
