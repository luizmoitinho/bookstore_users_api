package users

import (
	"fmt"

	"github.com/luizmoitinho/bookstore_users_api/date_utils"
	"github.com/luizmoitinho/bookstore_users_api/util/errors"
)

var usersDB = make(map[int64]*UserDTO)

func (user *UserDTO) Get() *errors.RestError {
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
	current := usersDB[user.Id]
	if current != nil {
		if current.Email == user.Email {
			return errors.NewBadRequestError(fmt.Sprintf("email %s already registered", user.Email))
		}
		return errors.NewBadRequestError(fmt.Sprintf("user %d already exists", user.Id))
	}
	user.CreatedAt = date_utils.GetNowString()
	return nil
}
