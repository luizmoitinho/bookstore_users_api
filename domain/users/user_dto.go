package users

import (
	"net/mail"
	"strings"

	"github.com/luizmoitinho/bookstore_users_api/util/errors"
)

const (
	STATUS_ACTIVE   = "active"
	STATUS_INACTIVE = "inactive"
)

type UserDTO struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	Password  string `json:"password"`
}

func (u *UserDTO) TreatmentAndValidate() *errors.RestError {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if u.Email == "" {
		return errors.NewBadRequestError("email address not be empty")
	} else if !ValidEmail(u.Email) {
		return errors.NewBadRequestError("email address is not valid")
	}

	if u.Password == "" {
		return errors.NewBadRequestError("password cannot be empty")
	}

	return nil
}

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
