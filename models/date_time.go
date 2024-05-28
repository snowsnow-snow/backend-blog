package models

import (
	"database/sql/driver"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

type DateTime time.Time

func (t *DateTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = DateTime(time.Time{})
		return
	}
	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = DateTime(now)
	return
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	if t.String() == "0001-01-01 00:00:00" {
		var empty []byte
		empty = append(empty, '"')
		empty = append(empty, '"')
		return empty, nil
	}
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t DateTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(TimeFormat)), nil
}

func (t *DateTime) Scan(v interface{}) error {
	date := string(v.([]uint8))
	if date == "0001-01-01 00:00:00" {
		return nil
	}
	tTime, _ := time.Parse(TimeFormat, string(v.([]uint8)))
	*t = DateTime(tTime)
	return nil
}

func (t DateTime) String() string {
	return time.Time(t).Format(TimeFormat)
}

// Now 获取当前时间
func (d DateTime) Now() DateTime {
	return DateTime(time.Now())
}
