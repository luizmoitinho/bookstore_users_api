package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luizmoitinho/bookstore_users_api/domain/users"
	"github.com/luizmoitinho/bookstore_users_api/services"
	"github.com/luizmoitinho/bookstore_users_api/util/errors"
)

func GetUser(c *gin.Context) {
	userId, err := getUserIdParams(c)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	result, getErr := services.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

func UpdateUser(c *gin.Context) {
	isPartial := c.Request.Method == http.MethodPatch

	userId, err := getUserIdParams(c)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	var user users.UserDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId
	if !isPartial {
		errTreatment := user.TreatmentAndValidate()
		if errTreatment != nil {
			c.JSON(errTreatment.Status, errTreatment)
			return
		}
		if user.FirstName == "" || user.LastName == "" {
			err := errors.NewBadRequestError("fist name or last name not be empty")
			c.JSON(err.Status, err)
			return
		}
	}

	result, saveErr := services.UpdateUser(isPartial, user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func getUserIdParams(c *gin.Context) (int64, *errors.RestError) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		return 0, errors.NewBadRequestError("user id should be a number")
	} else if userId <= 0 {
		return 0, errors.NewBadRequestError("user id should be more than zero")
	}
	return userId, nil
}

func CreateUser(c *gin.Context) {
	var user users.UserDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusCreated, result)
}
