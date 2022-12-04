package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luizmoitinho/bookstore_users_api/domain/users"
	"github.com/luizmoitinho/bookstore_users_api/services"
	"github.com/luizmoitinho/bookstore_users_api/util/date_utils"
	"github.com/luizmoitinho/bookstore_users_api/util/errors"
)

func Get(c *gin.Context) {
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

func Search(c *gin.Context) {
	status := c.Query("status")
	if strings.TrimSpace(status) == "" {
		err := errors.NewBadRequestError("status parameter is not valid")
		c.JSON(err.Status, err)
		return
	}

	users, err := services.Search(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, users)

}

func Delete(c *gin.Context) {
	userId, err := getUserIdParams(c)
	if err != nil {
		c.JSON(err.Status, err)
	}

	if errDelete := services.DeleteUser(userId); errDelete != nil {
		c.JSON(errDelete.Status, errDelete)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Update(c *gin.Context) {
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

	c.JSON(http.StatusOK, result)
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

func Create(c *gin.Context) {
	var user users.UserDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Status = users.STATUS_ACTIVE
	user.CreatedAt = date_utils.GetNowDbFormat()
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusCreated, result)
}
