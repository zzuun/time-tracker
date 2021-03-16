package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/zzuun/time-tracker/model"
	"github.com/zzuun/time-tracker/utils"
	"net/http"
	"strconv"
	"time"
)

type Controller struct{ DB *model.Database }

// @Tags account
// @Summary signup
// @Param body	body model.UserInput	true	"username, password"
// @Accept  json
// @Produce  json
// @router /signup [post]
// @Success 201 "user created successfully"
func (c *Controller) Signup(ctx *gin.Context) {
	user := new(model.User)
	err := ctx.ShouldBind(&user)
	if err != nil {
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

	_, err = c.DB.GetUser(user.Email)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, "email already registered")
		return
	}
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	hashpassword, _ := utils.HashPassword(user.Password)
	user.Password = hashpassword

	_, err = c.DB.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, "user created successfully")
}

// @Tags account
// @Summary login
// @Param body	body model.UserInput	true	"username, password"
// @Accept  json
// @Produce  json
// @router /login [post]
// @Success 200 "X-AUTH-TOKEN"
func (c *Controller) Login(ctx *gin.Context) {
	user := new(model.User)
	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	password := user.Password
	user, err = c.DB.GetUser(user.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, "user not found")
		return
	}

	if !utils.CheckPassword(password, user.Password) {
		ctx.JSON(http.StatusBadRequest, "incorrect password")
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, token)
}

// @Tags timer
// @Summary Start timer
// @param X-AUTH-TOKEN	header	string	true	"JWT Token"
// @Accept  json
// @Produce  json
// @router /start [post]
// @Success 200 {object} model.Entry
func (c *Controller) StartTime(ctx *gin.Context) {

	token := ctx.GetHeader("X-AUTH-TOKEN")
	if len(token) == 0 {
		ctx.JSON(http.StatusBadRequest, "token is missing")
		return
	}

	id, err := utils.ValidateToken(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	entry, err := c.DB.AddTimeEntry(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	entry.UserId = id
	ctx.JSON(http.StatusOK, entry)
}

// @Tags timer
// @Summary Stop timer
// @Param id path	string 	true 	"id of the entry"
// @param X-AUTH-TOKEN	header	string	true	"JWT Token"
// @Accept  json
// @Produce  json
// @router /stop/{id} [put]
// @Success 200 "record updated"
func (c *Controller) StopTime(ctx *gin.Context) {

	token := ctx.GetHeader("X-AUTH-TOKEN")
	if len(token) == 0 {
		ctx.JSON(http.StatusBadRequest, "token is missing")
		return
	}

	user_id, err := utils.ValidateToken(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	strid := ctx.Param("id")
	id, _ := strconv.Atoi(strid)
	err = c.DB.UpdateTimeEntry(id, user_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "record updated")
}

// @Tags timer
// @Summary List of all entries
// @param X-AUTH-TOKEN	header	string	true	"JWT Token"
// @param from	query	string	false	"starting date : format 2021-01-01"
// @param to	query	string	false	"ending date : format 2021-01-31"
// @Accept  json
// @Produce  json
// @router /activity [get]
// @Success 200 "{"today":"","24hours":"","weekly":"","monthly":""}"
func (c *Controller) Activity(ctx *gin.Context) {

	token := ctx.GetHeader("X-AUTH-TOKEN")
	if len(token) == 0 {
		ctx.JSON(http.StatusBadRequest, "token is missing")
		return
	}

	user_id, err := utils.ValidateToken(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	var resp = make(map[string]string)

	to_input := ctx.Query("to")
	from_input := ctx.Query("from")

	//custom
	if len(to_input) != 0 && len(from_input) != 0 {
		to, err := utils.StringtoTime(to_input)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		from, err := utils.StringtoTime(from_input)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		if from.After(to) {
			from, to = to, from
		}
		to = to.AddDate(0, 0, 1)
		s, err := c.DB.ListTimeEntries(from, to, user_id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		resp["custom"] = s
	}

	// today
	from := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	to := time.Now()
	s, err := c.DB.ListTimeEntries(from, to, user_id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	resp["today"] = s

	//last24hours
	from = time.Now().AddDate(0, 0, -1)
	to = time.Now()
	s, err = c.DB.ListTimeEntries(from, to, user_id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	resp["24hours"] = s

	//todo this week
	to = time.Now()
	offset := int(time.Monday - to.Weekday())
	if offset > 0 {
		offset = -6
	}
	from = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	s, err = c.DB.ListTimeEntries(from, to, user_id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	resp["weekly"] = s

	// his month
	to = time.Now()
	from = time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, time.Local)
	s, err = c.DB.ListTimeEntries(from, to, user_id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	resp["monthly"] = s

	ctx.JSON(http.StatusOK, resp)

}
