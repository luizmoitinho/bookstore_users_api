package services

import (
	"github.com/luizmoitinho/bookstore_users_api/domain/users"
	"github.com/luizmoitinho/bookstore_users_api/util/errors"
)

func GetUser(userId int64) (*users.UserDTO, *errors.RestError) {
	result := &users.UserDTO{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func CreateUser(user users.UserDTO) (*users.UserDTO, *errors.RestError) {
	if err := user.TreatmentAndValidate(); err != nil {
		return nil, err
	}

	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}
