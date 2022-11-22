package users

import (
	"fmt"
	"strings"

	"github.com/luizmoitinho/bookstore_users_api/datasources/users_db"
	"github.com/luizmoitinho/bookstore_users_api/date_utils"
	"github.com/luizmoitinho/bookstore_users_api/util/errors"
)

const (
	UNIQUE_EMAIL = "users.UC_user_email"

	QUERY_INSERT = "INSERT INTO users (first_name, last_name, email, created_at) VALUES (?,?,?,?);"
)

var usersDB = make(map[int64]*UserDTO)

func (user *UserDTO) Get() *errors.RestError {
	con := users_db.Connect()
	if err := con.Client.Ping(); err != nil {
		panic(err)
	}

	result := usersDB[user.Id]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
	}

	user.Id = result.Id
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.CreatedAt = result.CreatedAt

	return nil
}

func (user *UserDTO) Save() *errors.RestError {
	conn := users_db.Connect()

	stm, err := conn.Client.Prepare(QUERY_INSERT)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stm.Close()

	user.CreatedAt = date_utils.GetNowString()
	result, err := stm.Exec(user.FirstName, user.LastName, user.Email, user.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), UNIQUE_EMAIL) {
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
