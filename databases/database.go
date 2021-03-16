package databases

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/zzuun/time-tracker/config"
	"github.com/zzuun/time-tracker/models"
	"time"
)

type Database interface {
	Close()
	CreateUser(models.User) (int, error)
	GetUser(string) (models.User, error)
	AddTimeEntry(int) (models.Entry, error)
	UpdateTimeEntry(int, int) error
	ListTimeEntries(time.Time, time.Time, int) (string, error)
}

func ConnectDatabase() (database Database, err error) {
	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Env.HOST, config.Env.PORT, config.Env.DBUSER, config.Env.PASSWORD, config.Env.DBNAME)
	client, err := sql.Open("postgres", info)
	if err != nil {
		return nil, err
	}
	database = &Postgres{client: client}
	return database, nil
}
