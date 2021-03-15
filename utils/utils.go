package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strconv"
	"time"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func StringtoTime(input string) (time.Time, error) {
	if len(input) < 0 && len(input) > 10 {
		return time.Time{}, errors.New("invalid date format")
	}

	str_y := input[:4]
	str_m := input[5:7]
	str_d := input[8:]

	y, err := strconv.Atoi(str_y)
	if err != nil {
		return time.Time{}, err
	}

	m, err := strconv.Atoi(str_m)
	if err != nil {
		return time.Time{}, err
	}
	d, err := strconv.Atoi(str_d)
	if err != nil {
		return time.Time{}, err
	}

	//if m < 1 || m > 12 {
	//	return time.Time{}, errors.New("invalid month")
	//}
	//if d < 1 || d > 31 {
	//	return time.Time{}, errors.New("invalid day")
	//}

	tm := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
	return tm, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(enterdPassword, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(enterdPassword))
	if err != nil {
		return false
	}
	return true
}
