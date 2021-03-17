package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/zzuun/time-tracker/config"
	"github.com/zzuun/time-tracker/models"
	"time"
)

type DataStore interface {
	Close()
	CreateUser(models.User) (int, error)
	GetUser(string) (models.User, error)
	AddTimeEntry(int) (models.Entry, error)
	UpdateTimeEntry(int, int) error
	ListTimeEntries(time.Time, time.Time, int) (string, error)
}

func NewDataStore() (datastore DataStore, err error) {
	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Conf.Host, config.Conf.Port, config.Conf.DbUser, config.Conf.Password, config.Conf.DbName)
	client, err := sql.Open("postgres", info)
	if err != nil {
		return nil, err
	}
	datastore = &Database{client: client}
	return datastore, nil
}
