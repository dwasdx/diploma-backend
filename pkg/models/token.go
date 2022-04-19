package models

import (
	"errors"
	"github.com/asaskevich/govalidator"
)

const FCMTokenPlatformIOS = "ios"
const FCMTokenPlatformAndriod = "android"

type FCMToken struct {
	Token    string `json:"token" valid:"stringlength(10|255)"`
	Platform string `json:"platform" valid:"stringlength(2|10)"`
}

func (s *FCMToken) Validate() (bool, error) {
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		return false, err
	}

	if !govalidator.IsIn(s.Platform, FCMTokenPlatformIOS, FCMTokenPlatformAndriod) {
		return false, errors.New("platform value is wrong")
	}
	return true, nil
}
