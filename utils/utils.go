package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func StringtoTime(input string) (time.Time, error) {
	return time.Parse("2006-01-02", input)
}
func IsEmailValid(e string) bool {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func ExtractUserId(id interface{}, exists bool) (int, error) {

	if !exists {
		return 0, errors.New("something went wrong")
	}

	return strconv.Atoi(fmt.Sprintf("%v", id))
}
