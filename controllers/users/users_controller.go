package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luizmoitinho/bookstore_users_api/domain/users"
	"github.com/luizmoitinho/bookstore_users_api/logger"
	"github.com/luizmoitinho/bookstore_users_api/services"
	"github.com/luizmoitinho/bookstore_users_api/util/crypto_utils"
	"github.com/luizmoitinho/bookstore_users_api/util/date_utils"
	"github.com/luizmoitinho/bookstore_users_api/util/errors"
)

func Get(c *gin.Context) {
	userId, err := getUserIdParams(c)
	if err != nil {
		logger.Error("error during get user id params at get route", errors.NewError(err.Error))
		c.JSON(err.Status, err)
		return
	}

	result, getErr := services.UsersService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, result.Marshall(isPublicRequest(c)))
}

func Search(c *gin.Context) {
	status := c.Query("status")
	if strings.TrimSpace(status) == "" {
		logger.Error("status parameter is not valid", nil)
		err := errors.NewBadRequestError("status parameter is not valid")
		c.JSON(err.Status, err)
		return
	}

	users, err := services.UsersService.SearchUser(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, users.Marshall(isPublicRequest(c)))
}

func Delete(c *gin.Context) {
	userId, err := getUserIdParams(c)
	if err != nil {
		logger.Error("error during get user id params at delete route", errors.NewError(err.Error))
		c.JSON(err.Status, err)
	}

	if errDelete := services.UsersService.DeleteUser(userId); errDelete != nil {
		c.JSON(errDelete.Status, errDelete)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Update(c *gin.Context) {
	isPartial := c.Request.Method == http.MethodPatch

	userId, err := getUserIdParams(c)
	if err != nil {
		logger.Error("error during get user id params at update route", errors.NewError(err.Error))
		c.JSON(err.Status, err)
		return
	}

	var user users.UserDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error("invalid json body", err)
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId
	if !isPartial {
		errTreatment := user.TreatmentAndValidate()
		if errTreatment != nil {
			logger.Error("error when trying update user at treatment and validate user", errors.NewError(errTreatment.Error))
			c.JSON(errTreatment.Status, errTreatment)
			return
		}
		if user.FirstName == "" || user.LastName == "" {
			logger.Error("error when trying update user", errors.NewError("fist name or last name not be empty"))
			err := errors.NewBadRequestError("first name or last name not be empty")
			c.JSON(err.Status, err)
			return
		}
	}

	result, saveErr := services.UsersService.UpdateUser(isPartial, user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusOK, result.Marshall(isPublicRequest(c)))
}

func Create(c *gin.Context) {
	var user users.UserDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error("invalid json body", err)
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Status = users.STATUS_ACTIVE
	user.Password = crypto_utils.GetSha256(user.Password)
	user.CreatedAt = date_utils.GetNowDbFormat()
	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		logger.Error("error when creating a new user", errors.NewError(saveErr.Error))
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusCreated, result.Marshall(isPublicRequest(c)))
}

func getUserIdParams(c *gin.Context) (int64, *errors.RestError) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		logger.Error("user id should be a number", userErr)
		return 0, errors.NewBadRequestError("user id should be a number")
	} else if userId <= 0 {
		logger.Error("error when trying get user id params", errors.NewError("user id should be more than zero"))
		return 0, errors.NewBadRequestError("user id should be more than zero")
	}
	return userId, nil
}

func isPublicRequest(c *gin.Context) bool {
	return c.GetHeader("x-public") == "true"
}
