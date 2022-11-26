package mysql_utils

import (
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/luizmoitinho/bookstore_users_api/util/errors"
)

const (
	ERROR_NO_ROWS               = "no rows in result set"
	ERROR_DUPLICATED_USER_EMAIL = "users.UC_user_email"
	DUPLICATED_KEY              = 1062
)

func ParseError(err error) *errors.RestError {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ERROR_NO_ROWS) {
			return errors.NewNotFoundError("no registration match given id")
		}
		return errors.NewInternalServerError("error parsing database response")
	}

	switch sqlErr.Number {
	case DUPLICATED_KEY:
		if strings.Contains(err.Error(), ERROR_DUPLICATED_USER_EMAIL) {
			return errors.NewBadRequestError("email already exists")
		}
		return errors.NewBadRequestError("duplicated key")
	}
	return errors.NewInternalServerError("error at processing request")
}
