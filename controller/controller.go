package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/zzuun/time-tracker/auth"
	"github.com/zzuun/time-tracker/databases"
	"github.com/zzuun/time-tracker/models"
	"github.com/zzuun/time-tracker/utils"
	"net/http"
	"strconv"
	"time"
)

type Controller struct{ db databases.Database }

func NewController(db databases.Database) *Controller {
	return &Controller{db: db}
}

// @Tags account
// @Summary signup
// @Param body	body models.UserInput	true	"username, password"
// @Accept  json
// @Produce  json
// @router /signup [post]
// @Success 201 "user created successfully"
func (ctrl *Controller) SignupPost(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if !auth.IsEmailValid(user.Email) {
		ctx.JSON(http.StatusBadRequest, "invalid email")
		return
	}

	if len(user.Password) == 0 {
		ctx.JSON(http.StatusBadRequest, "password cannot be empty")
		return
	}

	_, err := ctrl.db.GetUser(user.Email)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, "email already registered")
		return
	}
	if err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	hashpassword, _ := auth.HashPassword(user.Password)
	user.Password = hashpassword

	if _, err := ctrl.db.CreateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, "user created successfully")
}

// @Tags account
// @Summary login
// @Param body	body models.UserInput	true	"username, password"
// @Accept  json
// @Produce  json
// @router /login [post]
// @Success 200 "X-Auth-Token"
func (ctrl *Controller) LoginPost(ctx *gin.Context) {

	var user models.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	password := user.Password
	user, err := ctrl.db.GetUser(user.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, "user not found")
		return
	}

	if !auth.CheckPassword(password, user.Password) {
		ctx.JSON(http.StatusBadRequest, "incorrect password")
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, token)
}

// @Tags timer
// @Summary Start timer
// @param X-Auth-Token	header	string	true	"JWT Token"
// @Accept  json
// @Produce  json
// @router /tracker/start [post]
// @Success 200 {object} models.Entry
func (ctrl *Controller) StartTimePost(ctx *gin.Context) {

	userIdString, _ := ctx.Get(utils.UserId)
	userId := userIdString.(int)

	entry, err := ctrl.db.AddTimeEntry(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	entry.UserId = userId
	ctx.JSON(http.StatusOK, entry)
}

// @Tags timer
// @Summary Stop timer
// @Param id path	string 	true 	"id of the entry"
// @param X-Auth-Token	header	string	true	"JWT Token"
// @Accept  json
// @Produce  json
// @router /tracker/stop/{id} [put]
// @Success 200 "record updated"
func (ctrl *Controller) StopTimePut(ctx *gin.Context) {

	userIdString, _ := ctx.Get(utils.UserId)
	userId := userIdString.(int)

	strid := ctx.Param("id")
	id, _ := strconv.Atoi(strid)
	if err := ctrl.db.UpdateTimeEntry(id, userId); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "record updated")
}

// @Tags timer
// @Summary List of all entries
// @param X-Auth-Token	header	string	true	"JWT Token"
// @param duration	query	string	true	"duration: custom,week,day,month,day,today "
// @param from	query	string	false	"starting date : format 2021-01-01"
// @param to	query	string	false	"ending date : format 2021-01-31"
// @Accept  json
// @Produce  json
// @router /tracker/activity [get]
// @Success 200 "{"total_time":""}"
func (ctrl *Controller) ActivityGet(ctx *gin.Context) {

	userIdString, _ := ctx.Get(utils.UserId)
	userId := userIdString.(int)

	var to, from time.Time
	var err error

	switch ctx.Query(utils.ActivityDuration) {
	case utils.ActivityCustom:

		toInput := ctx.Query(utils.ActivityParamTimeTo)
		fromInput := ctx.Query(utils.ActivityParamTimeFrom)
		if to, err = utils.StringtoTime(toInput); err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		if from, err = utils.StringtoTime(fromInput); err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

	case utils.ActivityDay:

		from = time.Now().AddDate(0, 0, -1)
		to = time.Now()

	case utils.ActivityMonthly:

		to = time.Now()
		from = time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, time.Local)

	case utils.ActivityWeekly:

		to = time.Now()
		offset := int(time.Monday - to.Weekday())
		if offset > 0 {
			offset = -6
		}
		from = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)

	case utils.ActivityToday:

		from = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
		to = time.Now()

	default:
		ctx.JSON(http.StatusBadRequest, "duration is not valid")
		return
	}

	var resp = make(map[string]string)

	totalTime, err := ctrl.db.ListTimeEntries(from, to, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	resp["total_time"] = totalTime

	ctx.JSON(http.StatusOK, resp)

}
