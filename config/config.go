package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Configuration struct {
	HOST     string
	DBUSER   string
	PASSWORD string
	DBNAME   string
	PORT     string
}

var Env Configuration

func init() {
	if err := godotenv.Load(".env"); err != nil {
		return
	}

	Env.DBUSER = os.Getenv("DBUSER")
	Env.HOST = os.Getenv("HOST")
	Env.PASSWORD = os.Getenv("PASSWORD")
	Env.DBNAME = os.Getenv("DBNAME")
	Env.PORT = os.Getenv("PORT")
}
