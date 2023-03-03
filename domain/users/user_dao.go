package users

import (
	"fmt"
	"strings"

	"github.com/luizmoitinho/bookstore_users_api/datasources/users_db"
	"github.com/luizmoitinho/bookstore_users_api/logger"
	"github.com/luizmoitinho/bookstore_users_api/util/errors"
	"github.com/luizmoitinho/bookstore_users_api/util/mysql_utils"
)

const (
	QUERY_INSERT_USER                = "INSERT INTO users (first_name, last_name, email, created_at, status, password) VALUES (?,?,?,?,?,?);"
	QUERY_GET_USER                   = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE id=?;"
	QUERY_UPDATE_USER                = "UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?;"
	QUERY_DELETE_USER                = "DELETE FROM users WHERE id = ?;"
	QUERY_FIND_BY_STATUS             = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE status = ?;"
	QUERY_FIND_BY_EMAIL_AND_PASSWORD = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE email=? AND password=? AND status=?;"
)

func (user *UserDTO) Get() *errors.RestError {
	conn := users_db.Connect()

	query, err := conn.Client.Query(QUERY_GET_USER, user.Id)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError(errors.DATABASE_ERROR)
	}
	defer query.Close()

	if query.Next() {
		if getErr := query.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.Status); getErr != nil {
			logger.Error("error when trying to scan user row into user struct", getErr)
			return errors.NewInternalServerError(errors.DATABASE_ERROR)
		}
		return nil
	}
	return errors.NewNotFoundError("user with that given id not found")
}

func (user *UserDTO) FindByEmailAndPassword() *errors.RestError {
	conn := users_db.Connect()

	stm, err := conn.Client.Prepare(QUERY_FIND_BY_EMAIL_AND_PASSWORD)
	if err != nil {
		logger.Error("error when trying to prepare find user by email and password statement", err)
		return errors.NewInternalServerError(errors.DATABASE_ERROR)
	}
	defer stm.Close()

	result := stm.QueryRow(user.Email, user.Password, STATUS_ACTIVE)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ERROR_NO_ROWS) {
			return errors.NewNotFoundError("no user found with given credentials")
		}
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

	stm, err := conn.Client.Prepare(QUERY_FIND_BY_STATUS)
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
