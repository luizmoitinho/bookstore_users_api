package users

import (
	"github.com/luizmoitinho/bookstore_users_api/datasources/users_db"
	"github.com/luizmoitinho/bookstore_users_api/util/date_utils"
	"github.com/luizmoitinho/bookstore_users_api/util/errors"
	"github.com/luizmoitinho/bookstore_users_api/util/mysql_utils"
)

const (
	QUERY_INSERT_USER = "INSERT INTO users (first_name, last_name, email, created_at) VALUES (?,?,?,?);"
	QUERY_GET_USER    = "SELECT id, first_name, last_name, email, created_at FROM users WHERE id=?;"
)

func (user *UserDTO) Get() *errors.RestError {
	conn := users_db.Connect()

	stm, err := conn.Client.Prepare(QUERY_GET_USER)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stm.Close()

	result := stm.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt); getErr != nil {
		return mysql_utils.ParseError(getErr)
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
	result, saveErr := stm.Exec(user.FirstName, user.LastName, user.Email, user.CreatedAt)
	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	user.Id = userId

	return nil
}
