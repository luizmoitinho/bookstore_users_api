package services

import (
	"github.com/luizmoitinho/bookstore_users_api/domain/users"
	"github.com/luizmoitinho/bookstore_utils/rest_errors"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct{}

type usersServiceInterface interface {
	GetUser(int64) (*users.UserDTO, *rest_errors.RestError)
	SearchUser(string) (users.Users, *rest_errors.RestError)
	DeleteUser(int64) *rest_errors.RestError
	UpdateUser(bool, users.UserDTO) (*users.UserDTO, *rest_errors.RestError)
	CreateUser(users.UserDTO) (*users.UserDTO, *rest_errors.RestError)
	Authenticate(users.Login) (*users.UserDTO, *rest_errors.RestError)
}

func (s *usersService) GetUser(userId int64) (*users.UserDTO, *rest_errors.RestError) {
	result := &users.UserDTO{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *usersService) SearchUser(status string) (users.Users, *rest_errors.RestError) {
	user := users.UserDTO{}
	return user.FindByStatus(status)
}

func (s *usersService) Authenticate(login users.Login) (*users.UserDTO, *rest_errors.RestError) {
	user := users.UserDTO{Email: login.Email, Password: login.Password}
	if err := user.TreatmentAndValidate(); err != nil {
		return nil, err
	}

	if err := user.FindByEmailAndPassword(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) DeleteUser(userId int64) *rest_errors.RestError {
	user := &users.UserDTO{Id: userId}
	return user.Delete()
}

func (s *usersService) UpdateUser(isPartial bool, user users.UserDTO) (*users.UserDTO, *rest_errors.RestError) {
	currentUser, err := s.GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	if err := currentUser.TreatmentAndValidate(); err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			currentUser.FirstName = user.FirstName
		}
		if user.LastName != "" {
			currentUser.LastName = user.LastName
		}
		if user.Email != "" {
			currentUser.Email = user.Email
		}
	} else {
		currentUser.FirstName = user.FirstName
		currentUser.LastName = user.LastName
		currentUser.Email = user.Email
	}

	if err := currentUser.Update(); err != nil {
		return nil, err
	}

	return currentUser, nil
}

func (s *usersService) CreateUser(user users.UserDTO) (*users.UserDTO, *rest_errors.RestError) {
	if err := user.TreatmentAndValidate(); err != nil {
		return nil, err
	}

	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}
