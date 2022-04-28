package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// Claims - custom jwt token claims
type Claims struct {
	UserUUID string
	*jwt.StandardClaims
}

// Valid claims
func (c Claims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}
	if c.UserUUID == "" {
		return errors.New("bad user id")
	}
	return nil
}
