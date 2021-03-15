package model

import "time"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" binding:"Required"`
	Password string `json:"password,omitempty" binding:"Required"`
}

type Entry struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	StartTime time.Time `json:"start_time"`
	//EndTime   time.Time `json:"end_time,omitempty"`
	//Date      time.Time
}
