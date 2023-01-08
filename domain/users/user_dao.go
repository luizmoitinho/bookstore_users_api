package users

import (
	"fmt"

	"github.com/luizmoitinho/bookstore_users_api/datasources/users_db"
	"github.com/luizmoitinho/bookstore_users_api/logger"
	"github.com/luizmoitinho/bookstore_users_api/util/errors"
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
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError(errors.DATABASE_ERROR)
	}
	defer stm.Close()

	result := stm.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.Status); getErr != nil {
		logger.Error("error when trying to scan user row into user struct", getErr)
		return errors.NewInternalServerError(errors.DATABASE_ERROR)
	}

	return nil
}

func (user *UserDTO) Update() *errors.RestError {
	conn := users_db.Connect()

	stm, err := conn.Client.Prepare(QUERY_UPDATE_USER)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServerError(errors.DATABASE_ERROR)
	}
	defer stm.Close()

	_, updatedErr := stm.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if updatedErr != nil {
		logger.Error("error when trying update user", updatedErr)
		return errors.NewInternalServerError(errors.DATABASE_ERROR)
	}

	return nil
}

func (user *UserDTO) Delete() *errors.RestError {
	conn := users_db.Connect()

	stm, err := conn.Client.Prepare(QUERY_DELETE_USER)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError(errors.DATABASE_ERROR)
	}
	defer stm.Close()

	_, errDelete := stm.Exec(user.Id)
	if errDelete != nil {
		logger.Error("error when trying to delete user", err)
		return errors.NewInternalServerError(errors.DATABASE_ERROR)
	}
	return nil
}

func (user *UserDTO) Save() *errors.RestError {
	conn := users_db.Connect()

	stm, err := conn.Client.Prepare(QUERY_INSERT_USER)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServerError(errors.DATABASE_ERROR)
	}
	defer stm.Close()

	result, saveErr := stm.Exec(user.FirstName, user.LastName, user.Email, user.CreatedAt, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return errors.NewInternalServerError(errors.DATABASE_ERROR)
	}

	userId, err := result.LastInsertId()
	if err != nil {
		logger.Error("error when trying get last insert id after creating a new user", err)
		return errors.NewInternalServerError(errors.DATABASE_ERROR)
	}
	user.Id = userId

	return nil
}

func (user *UserDTO) FindByStatus(status string) ([]UserDTO, *errors.RestError) {
	conn := users_db.Connect()

	stm, err := conn.Client.Prepare(QUERY_FIND_USER_BY_STATUS)
	if err != nil {
		logger.Error("error when trying to prepare find users by status statement", err)
		return nil, errors.NewInternalServerError(errors.DATABASE_ERROR)
	}
	defer stm.Close()

	rows, err := stm.Query(status)
	if err != nil {
		logger.Error("error when trying to find users by status", err)
		return nil, errors.NewInternalServerError(errors.DATABASE_ERROR)
	}
	defer rows.Close()

	result := make([]UserDTO, 0)
	for rows.Next() {
		var user UserDTO
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.Status); err != nil {
			logger.Error("error when trying to scan user row into user struct", err)
			return nil, errors.NewInternalServerError(errors.DATABASE_ERROR)
		}
		result = append(result, user)
	}

	if len(result) == 0 {
		logger.Warn(fmt.Sprintf("no users matching status %s", status))
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return result, nil
}
