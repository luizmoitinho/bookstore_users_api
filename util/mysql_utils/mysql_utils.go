package mysql_utils

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/luizmoitinho/bookstore_users_api/util/errors"
)

const (
	ERROR_NO_ROWS  = "no rows in result set"
	DUPLICATED_KEY = 1062
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
		return errors.NewBadRequestError(fmt.Sprintf("duplicated key: %v", sqlErr.Message))
	}
	return errors.NewInternalServerError("error at processing request")
}
