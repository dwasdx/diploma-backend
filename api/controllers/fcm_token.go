package controllers

import (
	"encoding/json"
	"net/http"
	"shopingList/api"
	"shopingList/api/auth"
	"shopingList/pkg/models"
	"shopingList/store"
)

type FCMTokenController struct {
	authService *auth.Service
	storage     store.FCMTokenStorage
}

func NewFCMTokenController(authService *auth.Service, storage store.FCMTokenStorage) *FCMTokenController {
	return &FCMTokenController{authService: authService, storage: storage}
}

func (s *FCMTokenController) saveToken(w http.ResponseWriter, r *http.Request) {
	user, err := GetAuthorizedUser(s.authService, r)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusUnauthorized, err, "can't authorized user", api.ErrUserNotFound)
		return
	}

	tokenForm := models.FCMToken{}

	err = json.NewDecoder(r.Body).Decode(&tokenForm)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "error decode request", api.ErrValidationData)
		return
	}

	_, err = tokenForm.Validate()
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "validation error", api.ErrValidationData)
		return
	}

	result, err := s.storage.IsExistUserToken(user.ID, tokenForm.Token)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "error search exist token", api.ErrInternal)
		return
	}

	if !result {
		err = s.storage.SaveToken(user.ID, tokenForm)
		if err != nil {
			api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "error save token", api.ErrInternal)
			return
		}
	}

	api.SendDataJSON(w, r, http.StatusOK, nil)
}

func (s *FCMTokenController) deleteToken(w http.ResponseWriter, r *http.Request) {
	user, err := GetAuthorizedUser(s.authService, r)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusUnauthorized, err, "can't authorized user", api.ErrUserNotFound)
		return
	}

	tokenForm := models.FCMToken{}

	err = json.NewDecoder(r.Body).Decode(&tokenForm)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "error decode request", api.ErrValidationData)
		return
	}

	_, err = tokenForm.Validate()
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "validation error", api.ErrValidationData)
		return
	}

	result, err := s.storage.IsExistUserToken(user.ID, tokenForm.Token)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "error search exist token", api.ErrInternal)
		return
	}

	if !result {
		api.SendErrorJSON(w, r, http.StatusNotFound, err, "token not found", api.ErrValidationData)
		return
	}

	err = s.storage.DeleteUserToken(user.ID, tokenForm.Token)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "error save token", api.ErrInternal)
		return
	}

	api.SendDataJSON(w, r, http.StatusOK, nil)
}

func (s *FCMTokenController) Routes() []api.Route {
	return []api.Route{
		{
			Name:   "SaveToken",
			Method: "POST",
			Path:   "/fcm/token",
			Func:   s.saveToken,
		},
		{
			Name:   "DeleteToken",
			Method: "DELETE",
			Path:   "/fcm/token",
			Func:   s.deleteToken,
		},
	}
}
