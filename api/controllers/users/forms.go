package users

import (
	"errors"
	"github.com/asaskevich/govalidator"
)

type LoginForm struct {
	Phone int64 `json:"phone" valid:"numeric,stringlength(11|15),required"`
}

func (s *LoginForm) Validate() (bool, error) {
	if s.Phone == 0 {
		return false, errors.New("phone is empty")
	}

	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		return false, err
	}

	return true, nil
}

type PhoneConfirm struct {
	UserID string `json:"user_id" valid:"uuid,required"`
	Code   string `json:"code" valid:"numeric,required,stringlength(4|10)"`
}

func (s *PhoneConfirm) Validate() (bool, error) {
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		return false, err
	}

	if s.Code == "" {
		return false, errors.New("code is empty")
	}

	if s.UserID == "" {
		return false, errors.New("userId is empty")
	}

	return true, nil
}
