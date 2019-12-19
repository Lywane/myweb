package mysql

import (
	"strconv"
	"time"
	"fmt"
)

func ToString(param interface{}) string {
	switch ret := param.(type) {
	case string:
		return ret
	case []byte:
		return string(ret)
	case int:
		return strconv.Itoa(ret)
	case int64:
		return strconv.FormatInt(ret, 10)
	case float64:
		return strconv.FormatFloat(ret, 'f', -1, 64)
	case bool:
		if ret {
			return "1"
		} else {
			return "0"
		}
	case time.Time:
		return fmt.Sprint(ret)
	case nil:
		return ""
	default:
		return ""
	}
}

// 类型转换，任何类型转成int
func ToInt(param interface{}) int {
	switch ret := param.(type) {
	case int:
		return ret
	case int64:
		return int(ret)
	case float64:
		return int(ret)
	case []byte:
		r, _ := strconv.Atoi(string(ret))
		return r
	case string:
		r, _ := strconv.Atoi(ret)
		return r
	case bool:
		if ret {
			return 1
		} else {
			return 0
		}
	case nil:
		return 0
	default:
		return 0
	}
}

// 类型转换，任何类型转成int64
func ToInt64(param interface{}) int64 {
	switch ret := param.(type) {
	case int:
		return int64(ret)
	case int64:
		return ret
	case float64:
		return int64(ret)
	case []byte:
		r, _ := strconv.ParseInt(string(ret), 10, 64)
		return r
	case string:
		r, _ := strconv.ParseInt(ret, 10, 64)
		return r
	case bool:
		if ret {
			return 1
		} else {
			return 0
		}
	case time.Time:
		return ret.UnixNano() / 1000000
	case nil:
		return 0
	default:
		return 0
	}
}

// 类型转换，类型转换成float
func ToFloat(param interface{}) float64 {
	switch ret := param.(type) {
	case int64:
		return float64(ret)
	case float32:
		return float64(ret)
	case float64:
		return ret
	case []byte:
		r, _ := strconv.ParseFloat(string(ret), 64)
		return r
	case string:
		r, _ := strconv.ParseFloat(ret, 64)
		return r
	case bool:
		if ret {
			return 1.0
		} else {
			return 0.0
		}
	case nil:
		return 0.0
	default:
		return 0.0
	}
}

// 类型转换，任何类型转成bool
func ToBool(param interface{}) bool {
	switch ret := param.(type) {
	case bool:
		return ret
	case int64:
		if ret > 0 {
			return true
		} else {
			return false
		}
	case float64:
		if ret > 0 {
			return true
		} else {
			return false
		}
	case []byte:
		switch string(ret) {
		case "1", "true", "y", "on", "yes":
			return true
		case "0", "false", "n", "off", "no":
			return false
		default:
		}
		return false
	case string:
		switch ret {
		case "1", "true", "y", "on", "yes":
			return true
		case "0", "false", "n", "off", "no":
			return false
		default:
		}
		return false
	case nil:
		return false
	default:
		return false
	}
}

func ToTime(param interface{}) time.Time {
	switch ret := param.(type) {
	case []byte:
		r, err := time.ParseInLocation("2006-01-02 15:04:05", string(ret), time.Now().Location())
		if err != nil {
			return time.Now()
		}
		return r
	case string:
		r, err := time.ParseInLocation("2006-01-02 15:04:05", ret, time.Now().Location())
		if err != nil {
			return time.Now()
		}
		return r
	case time.Time:
		return ret
	default:
		return time.Now()
	}
}

