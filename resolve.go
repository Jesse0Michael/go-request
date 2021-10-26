package request

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// resolve the string value to the proper type and return the value
func resolve(t interface{}, v string) (interface{}, error) {
	switch t.(type) {
	case string:
		return v, nil
	case bool:
		return strconv.ParseBool(v)
	case time.Time:
		return time.Parse(time.RFC3339, v)
	case time.Duration:
		return time.ParseDuration(v)
	case int:
		i, err := strconv.ParseInt(v, 10, 32)
		return int(i), err
	case int64:
		return strconv.ParseInt(v, 10, 64)
	case int32:
		i, err := strconv.ParseInt(v, 10, 32)
		return int32(i), err
	case int16:
		i, err := strconv.ParseInt(v, 10, 16)
		return int16(i), err
	case int8:
		i, err := strconv.ParseInt(v, 10, 8)
		return int8(i), err
	case float64:
		return strconv.ParseFloat(v, 64)
	case float32:
		i, err := strconv.ParseFloat(v, 32)
		return float32(i), err
	case uint:
		i, err := strconv.ParseUint(v, 10, 32)
		return uint(i), err
	case uint64:
		return strconv.ParseUint(v, 10, 64)
	case uint32:
		i, err := strconv.ParseUint(v, 10, 32)
		return uint32(i), err
	case uint16:
		i, err := strconv.ParseUint(v, 10, 16)
		return uint16(i), err
	case uint8:
		i, err := strconv.ParseUint(v, 10, 8)
		return uint8(i), err
	case complex128:
		return strconv.ParseComplex(v, 128)
	case complex64:
		i, err := strconv.ParseComplex(v, 64)
		return complex64(i), err
	default:
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(t))
	}
}
