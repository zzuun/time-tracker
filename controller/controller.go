package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/zzuun/time-tracker/auth"
	"github.com/zzuun/time-tracker/db"
	"github.com/zzuun/time-tracker/models"
	"github.com/zzuun/time-tracker/utils"
	"net/http"
	"strconv"
	"time"
)

const (
	entryIDParam          = "id"
	activityParamTimeTo   = "to"
	activityParamTimeFrom = "from"
	activityCustom        = "custom"
	activityToday         = "today"
	activityWeekly        = "week"
	activityMonthly       = "month"
	activityDay           = "day"
	activityDuration      = "duration"
)

type Controller struct{ ds db.DataStore }

func NewController(ds db.DataStore) *Controller {
	return &Controller{ds: ds}
}

// @Tags account
// @Summary signup
// @Param body	body models.UserInput	true	"username, password"
// @Accept  json
// @Produce  json
// @router /signup [post]
// @Success 201 "user created successfully"
func (ctrl *Controller) SignupPOST(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if !utils.IsEmailValid(user.Email) {
		ctx.JSON(http.StatusBadRequest, "invalid email")
		return
	}

	if len(user.Password) == 0 {
		ctx.JSON(http.StatusBadRequest, "password cannot be empty")
		return
	}

	_, err := ctrl.ds.GetUser(user.Email)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, "email already registered")
		return
	}
	if err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	hashPassword, _ := auth.HashPassword(user.Password)
	user.Password = hashPassword

	if _, err := ctrl.ds.CreateUser(user); err != nil {
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
func (ctrl *Controller) LoginPOST(ctx *gin.Context) {

	var user models.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	password := user.Password
	user, err := ctrl.ds.GetUser(user.Email)
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
func (ctrl *Controller) StartTimePOST(ctx *gin.Context) {

	userId, err := utils.ExtractUserId(ctx.Get(utils.UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	entry, err := ctrl.ds.AddTimeEntry(userId)
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
// @router /tracker/stop/entry/{id} [put]
// @Success 200 "record updated"
func (ctrl *Controller) StopTimePUT(ctx *gin.Context) {

	userId, err := utils.ExtractUserId(ctx.Get(utils.UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	strid := ctx.Param(entryIDParam)
	id, err := strconv.Atoi(strid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if err := ctrl.ds.UpdateTimeEntry(id, userId); err != nil {
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
func (ctrl *Controller) ActivityGET(ctx *gin.Context) {

	userId, err := utils.ExtractUserId(ctx.Get(utils.UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	var to, from time.Time

	switch ctx.Query(activityDuration) {
	case activityCustom:

		toInput := ctx.Query(activityParamTimeTo)
		fromInput := ctx.Query(activityParamTimeFrom)
		if to, err = utils.StringtoTime(toInput); err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		if from, err = utils.StringtoTime(fromInput); err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

	case activityDay:

		from = time.Now().AddDate(0, 0, -1)
		to = time.Now()

	case activityMonthly:

		to = time.Now()
		from = time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, time.Local)

	case activityWeekly:

		to = time.Now()
		offset := int(time.Monday - to.Weekday())
		if offset > 0 {
			offset = -6
		}
		from = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)

	case activityToday:

		from = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
		to = time.Now()

	default:
		ctx.JSON(http.StatusBadRequest, "duration is not valid")
		return
	}

	var resp = make(map[string]string)

	totalTime, err := ctrl.ds.ListTimeEntries(from, to, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	resp["total_time"] = totalTime

	ctx.JSON(http.StatusOK, resp)

}
