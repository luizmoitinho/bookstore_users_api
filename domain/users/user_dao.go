package users

import (
	"fmt"

	"github.com/luizmoitinho/bookstore_users_api/datasources/users_db"
	"github.com/luizmoitinho/bookstore_users_api/util/errors"
	"github.com/luizmoitinho/bookstore_users_api/util/mysql_utils"
)

const (
	QUERY_INSERT_USER         = "INSERT INTO users (first_name, last_name, email, created_at, status, password) VALUES (?,?,?,?,?,?);"
	QUERY_GET_USER            = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE id=?;"
	QUERY_UPDATE_USER         = "UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?;"
	QUERY_DELETE_USER         = "DELETE FROM users WHERE id = ?;"
	QUERY_FIND_USER_BY_STATUS = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE status = ?;"
)

func (user *UserDTO) Get() *errors.RestError {
	conn := users_db.Connect()

	stm, err := conn.Client.Prepare(QUERY_GET_USER)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stm.Close()

	result := stm.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.Status); getErr != nil {
		return mysql_utils.ParseError(getErr)
	}

	return nil
}

func (user *UserDTO) Update() *errors.RestError {
	conn := users_db.Connect()

	stm, err := conn.Client.Prepare(QUERY_UPDATE_USER)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stm.Close()

	_, updatedErr := stm.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if updatedErr != nil {
		return mysql_utils.ParseError(updatedErr)
	}

	return nil
}

func (user *UserDTO) Delete() *errors.RestError {
	conn := users_db.Connect()

	stm, err := conn.Client.Prepare(QUERY_DELETE_USER)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stm.Close()

	_, errDelete := stm.Exec(user.Id)
	if errDelete != nil {
		return mysql_utils.ParseError(errDelete)
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

	result, saveErr := stm.Exec(user.FirstName, user.LastName, user.Email, user.CreatedAt, user.Status, user.Password)
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

func (user *UserDTO) FindByStatus(status string) ([]UserDTO, *errors.RestError) {
	conn := users_db.Connect()

	stm, err := conn.Client.Prepare(QUERY_FIND_USER_BY_STATUS)
	if err != nil {
		return nil, mysql_utils.ParseError(err)
	}
	defer stm.Close()

	rows, err := stm.Query(status)
	if err != nil {
		return nil, mysql_utils.ParseError(err)
	}
	defer rows.Close()

	result := make([]UserDTO, 0)
	for rows.Next() {
		var user UserDTO
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.Status); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		result = append(result, user)
	}

	if len(result) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return result, nil
}
