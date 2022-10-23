package services

import (
	"github.com/luizmoitinho/bookstore_users_api/domain/users"
	"github.com/luizmoitinho/bookstore_users_api/util/errors"
)

func CreateUser(user users.User) (*users.User, *errors.RestError) {
	user.CreatedAt = "123123"
	return &user, &errors.RestError{}
}
