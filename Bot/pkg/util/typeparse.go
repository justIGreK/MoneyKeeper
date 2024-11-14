package util

import (
	"strings"
	"time"
)

const (
	TimeFormat     string = "15:04"
	DateYearformat     string = "2006-01-02"
	YearFormat string = "2006"
	Dateformat     string = "01-02"
	DateTimeformat string = "2006-01-02T15:04"
	
)

func IsDate(part string) bool {
	_, err := time.Parse(DateYearformat, part)
	if err == nil {
		return true
	}
	_, err = time.Parse(Dateformat, part)
	return err == nil
}

func IsTime(part string) bool {
	_, err := time.Parse(TimeFormat, part)
	return err == nil
}

func ParseDate(part string) (time.Time, error) {
	if len(part) == 5 {
		part = time.Now().Format(YearFormat)+"-"+part
	}
	return time.Parse(DateYearformat, part)
}

func ParseTime(part string, date time.Time) (time.Time, error) {
	parsedTime, err := time.Parse(TimeFormat, part)
	if err != nil {
		return time.Time{}, err
	}

	if date.IsZero() {
		date = time.Now()
	}

	return time.Date(date.Year(), date.Month(), date.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Local), nil
}

func ParseKeyValue(input string) (string, string, bool) {
	parts := strings.SplitN(input, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}
