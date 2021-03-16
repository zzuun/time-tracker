package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zzuun/time-tracker/controller"
	_ "github.com/zzuun/time-tracker/docs"
	"github.com/zzuun/time-tracker/model"
	"log"
)

// @title time-tracker
// @version 1.0
// @description A simple time-tracking application in golang.

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
	router.PUT("/stop/:id", ctrl.StopTime)
	router.GET("/activity", ctrl.Activity)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Fatal(router.Run(":8000"))

}
