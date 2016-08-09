package utils

import (
	"fmt"
	"strings"
	"time"
)

func ParseDate(s []string) (time.Time, error) {
	if len(s) < 2 {
		return time.Time{}, fmt.Errorf("Not enough strings passed")
	}
	splitDate := strings.Split(s[0], ".")
	splitHour := strings.Split(s[1], ":")

	if len(splitDate) != 3 {
		return time.Time{}, fmt.Errorf("There were not 3 fields in the date string")
	}
	year, month, day := splitDate[0], splitDate[1], splitDate[2]
	if len(splitHour) < 2 {
		return time.Time{}, fmt.Errorf("There were not 2 fields in the hour string")
	}
	hour, minute := splitHour[0], splitHour[1]
	d := fmt.Sprintf("%s %s %s %s %s", month, day, year, hour, minute)
	return time.Parse("01 02 2006 15 04", d)
}

func ParseDateString(s []string) (string, error) {
	t, err := ParseDate(s)
	if err != nil {
		return "", err
	}
	return t.String(), nil
}
