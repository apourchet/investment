package utils

import (
	"testing"
)

func TestParseDate(t *testing.T) {
	s1 := []string{"2006.01.02", "19:03"}
	_, err := ParseDate(s1)
	if err != nil {
		t.Fatal("ParseDate returned error for valid date")
	}
	s2 := []string{"2016.12.30", "19:03"}
	_, err = ParseDate(s2)
	if err != nil {
		t.Fatal("ParseDate returned error for valid date")
	}

	s3 := []string{"2016/01/02", "19:03"}
	_, err = ParseDate(s3)
	if err == nil {
		t.Fatal("ParseDate returned no error for invalid date")
	}

	s3 = []string{"2016 01 02", "19.03"}
	_, err = ParseDate(s3)
	if err == nil {
		t.Fatal("ParseDate returned no error for invalid time")
	}
	s3 = []string{"2016 01 02"}
	_, err = ParseDate(s3)
	if err == nil {
		t.Fatal("ParseDate returned no error for invalid date/time slice")
	}
}

func TestParseDateString(t *testing.T) {
	s2 := []string{"2016.01.02", "19:03"}
	ti, err := ParseDateString(s2)
	if err != nil {
		t.Fatal("Error stringifying a valid date")
	}
	if ti != "2016-01-02 19:03:00 +0000 UTC" {
		t.Fatal(ti)
		t.Fatal("Failed to stringify a valid date")
	}
}
