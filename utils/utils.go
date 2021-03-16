package utils

import (
	"time"
)

func StringtoTime(input string) (time.Time, error) {
	return time.Parse("2006-01-02", input)
}
