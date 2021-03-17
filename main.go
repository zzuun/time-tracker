package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zzuun/time-tracker/auth"
	"github.com/zzuun/time-tracker/controller"
	"github.com/zzuun/time-tracker/db"
	_ "github.com/zzuun/time-tracker/docs"
	"github.com/zzuun/time-tracker/utils"
	"log"
	"net/http"
)

const (
	XAuthToken = "X-Auth-Token"
)

// @title time-tracker
// @version 1.0
// @description A simple time-tracking application in golang.

func main() {

	router := gin.Default()

	ds, err := db.NewDataStore()
	if err != nil {
		log.Fatal(err)
	}
	defer ds.Close()

	ctrl := controller.NewController(ds)

	//routes
	router.POST("/signup", ctrl.SignupPOST)
	router.POST("/login", ctrl.LoginPOST)

	trackerRoutes := router.Group("/tracker")
	trackerRoutes.Use(Authenticate())
	trackerRoutes.POST("/start", ctrl.StartTimePOST)
	trackerRoutes.PUT("/stop/entry/:id", ctrl.StopTimePUT)
	trackerRoutes.GET("/activity", ctrl.ActivityGET)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Fatal(router.Run(":8000"))

}

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader(XAuthToken)
		if len(token) == 0 {
			ctx.JSON(http.StatusUnauthorized, "X-Auth-Token is missing")
			ctx.Abort()
		}

		userId, err := auth.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, err.Error())
			ctx.Abort()
		}
		ctx.Set(utils.UserID, userId)
		ctx.Next()
	}
}
