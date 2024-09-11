package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// TimeStamp 用于表示时间戳的自定义类型
type TimeStamp struct {
	Time time.Time
}

var dateTime = "2006-01-02 15:04:05"

// MarshalJSON 实现自定义的 JSON 编码，将时间转换为标准时间格式
func (t TimeStamp) MarshalJSON() ([]byte, error) {
	if t.Time.String() == "0001-01-01 00:00:00 +0000 UTC" {
		return json.Marshal("")
	}
	longtime := t.Time.UnixNano() / 1e6
	sec := longtime / 1000
	msec := longtime % 1000
	format := time.Unix(sec, msec*int64(time.Millisecond)).Format(dateTime)
	return json.Marshal(format)
}

// UnmarshalJSON 实现自定义的 JSON 解码，将标准时间格式解析为时间戳
func (t *TimeStamp) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	parsedTime, err := time.Parse(dateTime, str)
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}

// 实现 Valuer 接口
func (t TimeStamp) Value() (driver.Value, error) {
	if t.Time.String() == "0001-01-01 00:00:00 +0000 UTC" {
		return nil, nil
	}
	return t.Time.UnixNano() / 1e6, nil
}

// 实现 Scanner 接口
func (t *TimeStamp) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	millis := value.(int64)
	seconds := millis / 1000
	nanoseconds := (millis % 1000) * 1e6
	// 创建 time.Time 对象
	currTime := time.Unix(seconds, nanoseconds)
	// 将时间戳转换为 time.Time 对象
	*t = TimeStamp{currTime}
	return nil
}
func Now() TimeStamp {
	return TimeStamp{Time: time.Now()}
}
