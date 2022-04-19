package models

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"time"
)

// User model
type User struct {
	ID          string     `json:"id" valid:"uuid,required" db:"id"`
	Name        string     `json:"name" db:"name"`
	Phone       int64      `json:"phone" valid:"length(11|15),required" db:"phone"`
	Email       NullString `json:"email" db:"email"`
	Code        NullString `json:"-" db:"code"`
	CreatedAt   int64      `json:"created_at" valid:"int,required" db:"created_at"`
	UpdatedAt   int64      `json:"updated_at" valid:"int,required" db:"updated_at"`
	IsActivated bool       `json:"is_activated" valid:"optional" db:"is_activated"`
	IsDeleted   bool       `json:"is_deleted" valid:"-" db:"is_deleted"`
}

func (s *User) Validate() (bool, error) {
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		return false, err
	}

	if !s.Email.IsEmpty() && !govalidator.IsEmail(s.Email.String) {
		return false, errors.New("email format is wrong")
	}

	if !s.Code.IsEmpty() {
		if !govalidator.StringLength(s.Code.String, "4", "10") {
			return false, errors.New("code format is wrong")
		}
	}

	if len(s.Name) > 0 {
		if !govalidator.StringLength(s.Name, "1", "100") {
			return false, errors.New("code format is wrong")
		}
	}

	if s.CreatedAt == 0 {
		return false, errors.New("CreatedAt is 0")
	}

	if s.UpdatedAt == 0 {
		return false, errors.New("UpdatedAt is 0")
	}

	return true, nil
}

// ExternalUser model
type ExternalUser struct {
	ID          string `json:"id"`
	Phone       int64  `json:"phone"`
	IsActivated bool   `json:"is_activated"`
}

func CreateExternalUser(user *User) ExternalUser {
	extUser := ExternalUser{ID: user.ID, Phone: user.Phone, IsActivated: user.IsActivated}
	return extUser
}

type UserInterface interface {
	getPhone() int64
	getID() string
}

func (s User) getID() string {
	return s.ID
}

func (s User) getPhone() int64 {
	return s.Phone
}

func (s ExternalUser) getID() string {
	return s.ID
}

func (s ExternalUser) getPhone() int64 {
	return s.Phone
}

func (s *User) RegenerateAuthCode(generator UserAuthCodeGenerator) {
	s.Code.String = generator.GetCode()
	s.UpdatedAt = time.Now().Unix()
}

func (s *User) ClearAuthCode() {
	s.Code.String = ""
	s.UpdatedAt = time.Now().Unix()
}

type UserAuthCodeGenerator interface {
	GetCode() string
}
