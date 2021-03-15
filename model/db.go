package model

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

var (
	HOST     string
	USER     string
	PASSWORD string
)

type Database struct {
	db *sql.DB
}

func ConnectDB() (*Database, error) {

	info := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=postgres sslmode=disable", HOST, USER, PASSWORD)
	db, err := sql.Open("postgres", info)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil

}

func (db *Database) Close() {
	_ = db.db.Close()
}

func (db *Database) CreateUser(user *User) (int, error) {
	var id int
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	err := db.db.QueryRow(query, user.Email, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) GetUser(email string) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE email=$1`
	row := db.db.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	return &user, err
}

func (db *Database) AddTimeEntry(user_id int) (Entry, error) {
	var e Entry
	query := `INSERT INTO enteries (user_id, start_time) VALUES ($1, $2) RETURNING id,start_time`
	err := db.db.QueryRow(query, user_id, time.Now()).Scan(&e.ID, &e.StartTime)
	if err != nil {
		return Entry{}, err
	}
	return e, err
}

func (db *Database) UpdateTimeEntry(id, user_id int) error {
	query := `UPDATE enteries SET end_time=$1 WHERE user_id=$2 AND id=$3`
	row := db.db.QueryRow(query, time.Now(), user_id, id)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}

func (db *Database) ListTimeEntries(from, to time.Time, user_id int) (string, error) {

	query := `SELECT sum(age(start_time,end_time)) FROM enteries WHERE start_time>=$1 AND end_time<=$2 AND user_id=$3`
	row, err := db.db.Query(query, from, to, user_id)
	if err != nil {
		return "", err
	}
	defer row.Close()

	var t []uint8
	if row.Next() {
		err = row.Scan(&t)
		if err != nil {
			return "", err
		}
	}
	if t == nil {
		return "", nil
	}
	return string(t[1:9]), nil
}
