package forms

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"strings"
)

// RegForm - structure of registration form
type RegForm struct {
	Username string `json:"username" valid:"stringlength(2|100),required"`
	Email    string `json:"email"`
	Phone    int64  `json:"phone" valid:"numeric,stringlength(11|15),required"`
}

func (s *RegForm) Validate() (bool, error) {
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		return false, err
	}

	s.Username = strings.TrimLeft(s.Username, " ")
	s.Username = strings.TrimRight(s.Username, " ")

	if len(s.Username) == 0 {
		return false, errors.New("Username is empty")
	}

	if len(s.Email) > 0 && !govalidator.IsEmail(s.Email) {
		return false, errors.New("email is wrong")
	}

	return true, nil
}
