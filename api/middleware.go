package api

import (
	"net/http"
)

func (s *Rest) accessTokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, _, err := s.authenticator.JWTWithClaims(r)
		if err != nil {
			SendErrorJSON(w, r, http.StatusUnauthorized, err, "can't get token", ErrCredentials)
			return
		}
		if err := claims.Valid(); err != nil {
			SendErrorJSON(w, r, http.StatusUnauthorized, err, "invalid token", ErrInvalidToken)
			return
		}

		repository := s.dataService.GetUsersReadRepository()

		u, err := repository.GetUser(claims.UserUUID)
		if err != nil {
			SendErrorJSON(w, r, http.StatusUnauthorized, err, "can't get user", ErrUserNotFound)
			return
		}

		if u.IsDeleted {
			SendErrorJSON(w, r, http.StatusUnauthorized, err, "user has deleted", ErrUserHasDeleted)
			return
		}

		if !u.IsActivated {
			SendErrorJSON(w, r, http.StatusUnauthorized, err, "user is not activated", ErrUserIsNotActivated)
			return
		}

		r = s.authenticator.SetUserInfo(r, *u)

		next.ServeHTTP(w, r)
	})
}
