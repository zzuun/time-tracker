package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/zzuun/time-tracker/models"
	"time"
)

type Database struct {
	client *sql.DB
}

func (database *Database) Close() {
	_ = database.client.Close()
}

func (database *Database) CreateUser(user models.User) (id int, err error) {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	err = database.client.QueryRow(query, user.Email, user.Password).Scan(&id)
	return
}

func (database *Database) GetUser(email string) (user models.User, err error) {
	query := `SELECT * FROM users WHERE email=$1`
	err = database.client.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	return user, err
}

func (database *Database) AddTimeEntry(user_id int) (entry models.Entry, err error) {
	query := `INSERT INTO enteries (user_id, start_time) VALUES ($1, $2) RETURNING id,start_time`
	row := database.client.QueryRow(query, user_id, time.Now())
	err = row.Scan(&entry.ID, &entry.StartTime)
	return entry, err
}

func (database *Database) UpdateTimeEntry(entryId, userId int) (err error) {
	query := `UPDATE enteries SET end_time=$1 WHERE user_id=$2 AND id=$3`
	err = database.client.QueryRow(query, time.Now(), userId, entryId).Scan(nil)
	if err != sql.ErrNoRows {
		return err
	}
	return nil
}

func (database *Database) ListTimeEntries(timeFrom, timeTo time.Time, userId int) (string, error) {
	query := `SELECT sum(age(start_time,end_time)) FROM enteries WHERE start_time>=$1 AND end_time<=$2 AND user_id=$3`
	row, err := database.client.Query(query, timeFrom, timeTo, userId)
	if err != nil {
		return "", err
	}
	defer row.Close()

	var totalTime []uint8
	if row.Next() {
		if err := row.Scan(&totalTime); err != nil {
			return "", err
		}
	}
	if totalTime == nil {
		return "", nil
	}
	return string(totalTime[1:9]), nil
}
