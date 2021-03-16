package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zzuun/time-tracker/auth"
	"github.com/zzuun/time-tracker/controller"
	"github.com/zzuun/time-tracker/databases"
	_ "github.com/zzuun/time-tracker/docs"
	"github.com/zzuun/time-tracker/utils"
	"log"
	"net/http"
)

// @title time-tracker
// @version 1.0
// @description A simple time-tracking application in golang.

func main() {

	router := gin.Default()

	db, err := databases.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctrl := controller.NewController(db)

	//routes
	router.POST("/signup", ctrl.SignupPost)
	router.POST("/login", ctrl.LoginPost)

	trackerRoutes := router.Group("/tracker")
	trackerRoutes.Use(Authenticate())
	trackerRoutes.POST("/start", ctrl.StartTimePost)
	trackerRoutes.PUT("/stop/:id", ctrl.StopTimePut)
	trackerRoutes.GET("/activity", ctrl.ActivityGet)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Fatal(router.Run(":8000"))

}

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader(utils.XAuthToken)
		if len(token) == 0 {
			ctx.JSON(http.StatusUnauthorized, "X-Auth-Token is missing")
			ctx.Abort()
		}

		userId, err := auth.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, err.Error())
			ctx.Abort()
		}
		ctx.Set(utils.UserId, userId)
		ctx.Next()
	}
}
