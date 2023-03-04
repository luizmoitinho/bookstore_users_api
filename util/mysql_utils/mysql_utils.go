package mysql_utils

import (
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/luizmoitinho/bookstore_utils/rest_errors"
)

const (
	ERROR_NO_ROWS               = "no rows in result set"
	ERROR_DUPLICATED_USER_EMAIL = "users.UC_user_email"
	DUPLICATED_KEY              = 1062
)

func ParseError(err error) *rest_errors.RestError {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ERROR_NO_ROWS) {
			return rest_errors.NewNotFoundError("no registration match given id")
		}
		return rest_errors.NewInternalServerError("error parsing database response", err)
	}

	switch sqlErr.Number {
	case DUPLICATED_KEY:
		if strings.Contains(err.Error(), ERROR_DUPLICATED_USER_EMAIL) {
			return rest_errors.NewBadRequestError("email already exists")
		}
		return rest_errors.NewBadRequestError("duplicated key")
	}
	return rest_errors.NewInternalServerError("error at processing request", err)
}
