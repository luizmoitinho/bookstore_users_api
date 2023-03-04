package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luizmoitinho/bookstore_oauth/oauth"
	"github.com/luizmoitinho/bookstore_users_api/domain/users"
	"github.com/luizmoitinho/bookstore_users_api/logger"
	"github.com/luizmoitinho/bookstore_users_api/services"
	"github.com/luizmoitinho/bookstore_users_api/util/crypto_utils"
	"github.com/luizmoitinho/bookstore_users_api/util/date_utils"
	"github.com/luizmoitinho/bookstore_utils/rest_errors"
)

func Get(c *gin.Context) {
	if err := oauth.Authenticate(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}

	if callerId := oauth.GetCallerId(c.Request); callerId == 0 {
		err := rest_errors.NewUnauthorizedError("resource not available")
		c.JSON(err.Status, err)
		return
	}

	userId, err := getUserIdParams(c)
	if err != nil {
		logger.Error("error during get user id params at get route", rest_errors.NewError(err.Error))
		c.JSON(err.Status, err)
		return
	}

	user, getErr := services.UsersService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	if oauth.GetCallerId(c.Request) == user.Id {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

func Search(c *gin.Context) {
	status := c.Query("status")
	if strings.TrimSpace(status) == "" {
		logger.Error("status parameter is not valid", nil)
		err := rest_errors.NewBadRequestError("status parameter is not valid")
		c.JSON(err.Status, err)
		return
	}

	users, err := services.UsersService.SearchUser(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, users.Marshall(oauth.IsPublic(c.Request)))
}

func Delete(c *gin.Context) {
	userId, err := getUserIdParams(c)
	if err != nil {
		logger.Error("error during get user id params at delete route", rest_errors.NewError(err.Error))
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
		logger.Error("error during get user id params at update route", rest_errors.NewError(err.Error))
		c.JSON(err.Status, err)
		return
	}

	var user users.UserDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error("invalid json body", err)
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId
	if !isPartial {
		errTreatment := user.TreatmentAndValidate()
		if errTreatment != nil {
			logger.Error("error when trying update user at treatment and validate user", rest_errors.NewError(errTreatment.Error))
			c.JSON(errTreatment.Status, errTreatment)
			return
		}
		if user.FirstName == "" || user.LastName == "" {
			logger.Error("error when trying update user", rest_errors.NewError("fist name or last name not be empty"))
			err := rest_errors.NewBadRequestError("first name or last name not be empty")
			c.JSON(err.Status, err)
			return
		}
	}

	result, saveErr := services.UsersService.UpdateUser(isPartial, user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusOK, result.Marshall(oauth.IsPublic(c.Request)))
}

func Create(c *gin.Context) {
	var user users.UserDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error("invalid json body", err)
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Status = users.STATUS_ACTIVE
	user.Password = crypto_utils.GetSha256(user.Password)
	user.CreatedAt = date_utils.GetNowDbFormat()
	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		logger.Error("error when creating a new user", rest_errors.NewError(saveErr.Error))
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusCreated, result.Marshall(oauth.IsPublic(c.Request)))
}

func Authenticate(c *gin.Context) {
	var login users.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		restError := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restError.Status, restError)
		return
	}

	login.Password = crypto_utils.GetSha256(login.Password)
	user, err := services.UsersService.Authenticate(login)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusCreated, user.Marshall(oauth.IsPublic(c.Request)))
}

func getUserIdParams(c *gin.Context) (int64, *rest_errors.RestError) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		logger.Error("user id should be a number", userErr)
		return 0, rest_errors.NewBadRequestError("user id should be a number")
	} else if userId <= 0 {
		logger.Error("error when trying get user id params", rest_errors.NewError("user id should be more than zero"))
		return 0, rest_errors.NewBadRequestError("user id should be more than zero")
	}
	return userId, nil
}
