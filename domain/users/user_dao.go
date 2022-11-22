package users

import (
	"fmt"
	"strings"

	"github.com/luizmoitinho/bookstore_users_api/datasources/users_db"
	"github.com/luizmoitinho/bookstore_users_api/date_utils"
	"github.com/luizmoitinho/bookstore_users_api/util/errors"
)

const (
	NO_ROWS_RESULT_SET = "no rows in result set"
	UNIQUE_USER_EMAIL  = "users.UC_user_email"

	QUERY_INSERT_USER = "INSERT INTO users (first_name, last_name, email, created_at) VALUES (?,?,?,?);"
	QUERY_GET_USER    = "SELECT id, first_name, last_name, email, created_at FROM users WHERE id=?;"
)

var usersDB = make(map[int64]*UserDTO)

func (user *UserDTO) Get() *errors.RestError {
	conn := users_db.Connect()

	stm, err := conn.Client.Prepare(QUERY_GET_USER)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stm.Close()

	result := stm.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt); err != nil {
		if strings.Contains(err.Error(), NO_ROWS_RESULT_SET) {
			return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error trying get user %v: %v", user.Id, err.Error()))
	}

	return nil
}

func (user *UserDTO) Save() *errors.RestError {
	conn := users_db.Connect()

	stm, err := conn.Client.Prepare(QUERY_INSERT_USER)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stm.Close()

	user.CreatedAt = date_utils.GetNowString()
	result, err := stm.Exec(user.FirstName, user.LastName, user.Email, user.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), UNIQUE_USER_EMAIL) {
			return errors.NewBadRequestError(fmt.Sprintf("email %v already exists", user.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error whe trying save user: %v", err.Error()))
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to get last insert id: %v", err.Error()))
	}
	user.Id = userId

	return nil
}
