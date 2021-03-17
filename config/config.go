package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Configuration struct {
	Host     string
	DbUser   string
	Password string
	DbName   string
	Port     string
}

var Conf Configuration

func init() {

	if err := godotenv.Load(".env"); err != nil {
		log.Println(err)
	}

	Conf.DbUser = os.Getenv("DB_USER")
	Conf.Host = os.Getenv("HOST")
	Conf.Password = os.Getenv("PASSWORD")
	Conf.DbName = os.Getenv("DB_NAME")
	Conf.Port = os.Getenv("PORT")
}
