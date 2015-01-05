package form

import (
	"fmt"
	"reflect"
	"strconv"
)

// TODO return error
type converter func(interface{}) interface{}

var (
	boolType    = reflect.TypeOf(false)
	float32Type = reflect.TypeOf(float32(0))
	float64Type = reflect.TypeOf(float64(0))
	intType     = reflect.TypeOf(int(0))
	int8Type    = reflect.TypeOf(int8(0))
	int16Type   = reflect.TypeOf(int16(0))
	int32Type   = reflect.TypeOf(int32(0))
	int64Type   = reflect.TypeOf(int64(0))
	stringType  = reflect.TypeOf("")
	uintType    = reflect.TypeOf(uint(0))
	uint8Type   = reflect.TypeOf(uint8(0))
	uint16Type  = reflect.TypeOf(uint16(0))
	uint32Type  = reflect.TypeOf(uint32(0))
	uint64Type  = reflect.TypeOf(uint64(0))
)

// converters for basic types
var converters = map[reflect.Type]converter{
	boolType:    convertBool,
	float32Type: func(v interface{}) interface{} { return float32(convertFloat64(v, 32)) },
	float64Type: func(v interface{}) interface{} { return convertFloat64(v, 64) },
	intType:     func(v interface{}) interface{} { return int(convertInt64(v, 0)) },
	int8Type:    func(v interface{}) interface{} { return int8(convertInt64(v, 8)) },
	int16Type:   func(v interface{}) interface{} { return int16(convertInt64(v, 16)) },
	int32Type:   func(v interface{}) interface{} { return int32(convertInt64(v, 32)) },
	int64Type:   func(v interface{}) interface{} { return convertInt64(v, 64) },
	uintType:    func(v interface{}) interface{} { return uint(convertUint64(v, 0)) },
	uint8Type:   func(v interface{}) interface{} { return uint8(convertUint64(v, 8)) },
	uint16Type:  func(v interface{}) interface{} { return uint16(convertUint64(v, 16)) },
	uint32Type:  func(v interface{}) interface{} { return uint32(convertUint64(v, 32)) },
	uint64Type:  func(v interface{}) interface{} { return convertUint64(v, 64) },
	stringType:  func(v interface{}) interface{} { return fmt.Sprintf("%v", v) },
}

func convertBool(v interface{}) interface{} {
	switch v.(type) {
	case bool:
		return v.(bool)
	case int:
		return v.(int) != 0
	case int8:
		return v.(int8) != 0
	case int16:
		return v.(int16) != 0
	case int32:
		return v.(int32) != 0
	case int64:
		return v.(int64) != 0
	case uint:
		return v.(uint) != 0
	case uint8:
		return v.(uint8) != 0
	case uint16:
		return v.(uint16) != 0
	case uint32:
		return v.(uint32) != 0
	case uint64:
		return v.(uint64) != 0
	case float32:
		return v.(float32) != 0
	case float64:
		return v.(float64) != 0
	case string:
		var s = v.(string)
		if s == "yes" || s == "on" {
			return true
		} else if x, err := strconv.ParseBool(s); err == nil {
			return x
		} else if x, err := strconv.ParseFloat(s, 64); err == nil {
			return x != 0
		}
	}

	return nil
}

func convertFloat64(v interface{}, size int) float64 {
	switch v.(type) {
	case bool:
		if v.(bool) {
			return float64(1)
		}
		return float64(0)
	case int:
		return float64(v.(int))
	case int8:
		return float64(v.(int8))
	case int16:
		return float64(v.(int16))
	case int32:
		return float64(v.(int32))
	case int64:
		return float64(v.(int64))
	case uint:
		return float64(v.(uint))
	case uint8:
		return float64(v.(uint8))
	case uint16:
		return float64(v.(uint16))
	case uint32:
		return float64(v.(uint32))
	case uint64:
		return float64(v.(uint64))
	case float32:
		return float64(v.(float32))
	case float64:
		return v.(float64)
	case string:
		s := v.(string)
		if x, err := strconv.ParseFloat(s, size); err == nil {
			return x
		}
	}

	return 0
}

func convertInt64(v interface{}, size int) int64 {
	switch v.(type) {
	case bool:
		if v.(bool) {
			return int64(1)
		}
		return int64(0)
	case int:
		return int64(v.(int))
	case int8:
		return int64(v.(int8))
	case int16:
		return int64(v.(int16))
	case int32:
		return int64(v.(int32))
	case int64:
		return int64(v.(int64))
	case uint:
		return int64(v.(uint))
	case uint8:
		return int64(v.(uint8))
	case uint16:
		return int64(v.(uint16))
	case uint32:
		return int64(v.(uint32))
	case uint64:
		return int64(v.(uint64))
	case float32:
		return int64(v.(float32))
	case float64:
		return int64(v.(float64))
	case string:
		s := v.(string)
		if x, err := strconv.ParseInt(s, 10, size); err == nil {
			return x
		}
	}
	return 0
}

func convertUint64(v interface{}, size int) uint64 {
	switch v.(type) {
	case bool:
		if v.(bool) {
			return uint64(1)
		}
		return uint64(0)
	case int:
		return uint64(v.(int))
	case int8:
		return uint64(v.(int8))
	case int16:
		return uint64(v.(int16))
	case int32:
		return uint64(v.(int32))
	case int64:
		return uint64(v.(int64))
	case uint:
		return uint64(v.(uint))
	case uint8:
		return uint64(v.(uint8))
	case uint16:
		return uint64(v.(uint16))
	case uint32:
		return uint64(v.(uint32))
	case uint64:
		return uint64(v.(uint64))
	case float32:
		return uint64(v.(float32))
	case float64:
		return uint64(v.(float64))
	case string:
		s := v.(string)
		if x, err := strconv.ParseUint(s, 10, size); err == nil {
			return x
		}
	}

	return 0
}
