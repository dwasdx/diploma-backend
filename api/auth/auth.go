package auth

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"shopingList/internal"

	//"shopingList/api"
	"shopingList/pkg/models"
)

var defaultJWTHeaderKey = "Authorization"

type contextKey string

// Opts - auth service options stricture
type Opts struct {
	SigningKey []byte
	JWTHeader  string
}

// Service - auth service
type Service struct {
	Opts
}

// NewService returns new instance of auth service
func NewService(opts Opts) *Service {
	svc := &Service{Opts: opts}
	if svc.JWTHeader == "" {
		svc.JWTHeader = defaultJWTHeaderKey
	}
	return svc
}

// CreateJWT returns string of jwt
func (s *Service) CreateJWT(claims Claims) (string, error) {
	if claims.UserUUID == "" {
		return "", internal.WrapError("invalid uid", nil)
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := t.SignedString(s.SigningKey)
	if err != nil {
		return "", internal.WrapError("can's signing token", err)
	}

	return str, nil
}

// JWTWithClaims returns jwt token claims from request
func (s *Service) JWTWithClaims(r *http.Request) (Claims, string, error) {
	header := r.Header.Get(s.JWTHeader)
	if header == "" {
		return Claims{}, "", errors.New("header is empty")
	}

	tokenString := header[len("Bearer "):len(header)]

	if tokenString == "" {
		return Claims{}, "", errors.New("token not presented")
	}
	claims, err := s.Parse(tokenString)
	if err != nil {
		return Claims{}, "", err
	}
	return claims, tokenString, nil
}

// Parse claims from token string and returns its
func (s *Service) Parse(str string) (Claims, error) {
	parser := jwt.Parser{SkipClaimsValidation: true}
	t, err := parser.ParseWithClaims(str, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("signing method is not valid")
		}
		return s.SigningKey, nil
	})
	if err != nil {
		return Claims{}, internal.WrapError("can't parse token", err)
	}

	claims, ok := t.Claims.(*Claims)
	if !ok {
		return Claims{}, errors.New("invalid token claims type")
	}
	return *claims, nil
}

// SetUserInfo returns *http.Request with already set user model
func (s *Service) SetUserInfo(r *http.Request, user models.User) *http.Request { // nolint gocritic
	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKey("user"), user)
	return r.WithContext(ctx)
}

// UserInfo returns user model from *http.Request
func (s *Service) UserInfo(r *http.Request) (*models.User, error) {
	u, ok := r.Context().Value(contextKey("user")).(models.User)
	if !ok {
		return nil, errors.New("can't get user from request context")
	}
	return &u, nil
}
