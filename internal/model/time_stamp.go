package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TimeStamp 用于表示时间戳的自定义类型（DB里存毫秒；JSON里输出 yyyy-MM-dd HH:mm:ss）
type TimeStamp struct {
	Time time.Time
}

const dateTimeLayout = "2006-01-02 15:04:05"

func (t TimeStamp) IsZero() bool {
	return t.Time.IsZero()
}

// MarshalJSON：输出标准时间字符串；零值输出空串
func (t TimeStamp) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(t.Time.Format(dateTimeLayout))
}

// UnmarshalJSON：支持 "" / null / "yyyy-MM-dd HH:mm:ss" / 数字时间戳(秒/毫秒)
func (t *TimeStamp) UnmarshalJSON(data []byte) error {
	if t == nil {
		return fmt.Errorf("TimeStamp.UnmarshalJSON: receiver is nil")
	}

	// 处理 null
	if strings.TrimSpace(string(data)) == "null" {
		t.Time = time.Time{}
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		str = strings.TrimSpace(str)
		if str == "" {
			t.Time = time.Time{}
			return nil
		}

		// 纯数字：按秒/毫秒解析
		if isDigits(str) {
			v, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return err
			}
			t.Time = unixAuto(v)
			return nil
		}

		tt, err := time.ParseInLocation(dateTimeLayout, str, time.Local)
		if err != nil {
			return err
		}
		t.Time = tt
		return nil
	}

	// 如果不是 string，再尝试数字
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		t.Time = unixAuto(num)
		return nil
	}

	return fmt.Errorf("TimeStamp.UnmarshalJSON: invalid json: %s", string(data))
}

// Value：写入 DB（毫秒 int64）；零值写 NULL
func (t TimeStamp) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time.UnixMilli(), nil
}

// Scan：从 DB 读出（兼容 int64/int/float64/[]byte/string/time.Time）
func (t *TimeStamp) Scan(value any) error {
	if t == nil {
		return fmt.Errorf("TimeStamp.Scan: receiver is nil")
	}
	if value == nil {
		t.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return nil

	case int64:
		t.Time = unixAuto(v)
		return nil

	case int:
		t.Time = unixAuto(int64(v))
		return nil

	case float64:
		t.Time = unixAuto(int64(v))
		return nil

	case []byte:
		return t.scanString(string(v))

	case string:
		return t.scanString(v)

	default:
		return fmt.Errorf("TimeStamp.Scan: unsupported type %T (%v)", value, value)
	}
}

func (t *TimeStamp) scanString(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		t.Time = time.Time{}
		return nil
	}

	// 数字：时间戳
	if isDigits(s) {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		t.Time = unixAuto(v)
		return nil
	}

	// 文本时间
	tt, err := time.ParseInLocation(dateTimeLayout, s, time.Local)
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

// unixAuto：按数量级自动判断秒/毫秒（你 DB 存毫秒，这里也兼容秒）
func unixAuto(v int64) time.Time {
	// 大于 1e12 基本就是毫秒（2025年的秒级时间戳 ~ 1.7e9）
	if v > 1e12 {
		return time.UnixMilli(v)
	}
	return time.Unix(v, 0)
}

func isDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// Now 获取当前时间
func Now() TimeStamp {
	return TimeStamp{Time: time.Now()}
}
