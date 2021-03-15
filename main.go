package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zzuun/time-tracker/controller"
	"github.com/zzuun/time-tracker/model"
	"log"
)

func main() {

	//todo change to env variables
	model.HOST = "172.17.0.2"
	model.PASSWORD = "postgres"
	model.USER = "postgres"

	router := gin.Default()

	db, err := model.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctrl := &controller.Controller{DB: db}

	//routes
	router.POST("/signup", ctrl.Signup)
	router.POST("/login", ctrl.Login)
	router.POST("/start", ctrl.StartTime)
	router.PUT("/stop/	:id", ctrl.StopTime)
	router.GET("/activity", ctrl.Activity)

	log.Fatal(router.Run(":8000"))

}
