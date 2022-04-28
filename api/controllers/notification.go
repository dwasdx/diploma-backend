package controllers

import (
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"net/http"
	"shopingList/api"
	"shopingList/api/auth"
	"shopingList/pkg/models"
	"shopingList/pkg/readModels"
	"shopingList/pkg/repositories"
	"shopingList/pkg/services"
	"strconv"
)

type NotificationController struct {
	authService    *auth.Service
	repository     repositories.NotificationsRepository
	readRepository readModels.NotificationsReadRepository
	PushChannel    chan services.PushNotificationMessage
}

func NewNotificationController(
	authService *auth.Service,
	repository repositories.NotificationsRepository,
	readRepository readModels.NotificationsReadRepository) *NotificationController {
	return &NotificationController{authService: authService, repository: repository, readRepository: readRepository}
}

type NotificationPushForm struct {
	ID string `json:"id" valid:"uuid,required"`
}

func (s *NotificationPushForm) Validate() (bool, error) {
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		return false, err
	}

	return true, nil
}

type NotificationBatchForm struct {
	Page int `json:"page"`
}

func (s *NotificationBatchForm) Validate() (bool, error) {
	if s.Page < 0 {
		return false, errors.New("page is wrong")
	}

	return true, nil
}

// Routes returns slice of server routes
func (s *NotificationController) Routes() []api.Route {
	return []api.Route{
		{
			Name:   "Notification",
			Method: "GET",
			Path:   "/notifications/page/{page}",
			Func:   s.getNotifications,
		},
		{
			Name:   "Notification",
			Method: "GET",
			Path:   "/notifications/",
			Func:   s.getNotifications,
		},
		{
			Name:   "NotificationPush",
			Method: "POST",
			Path:   "/notification/push-to-token",
			Func:   s.notificationPush,
		},
	}
}

type NotificationBatch struct {
	Total int                    `json:"total"`
	Items *[]models.Notification `json:"items"`
}

func (s *NotificationController) getNotifications(w http.ResponseWriter, r *http.Request) {
	currentUser, err := GetAuthorizedUser(s.authService, r)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusUnauthorized, err, "can't auth user by token", api.ErrUserNotFound)
		return
	}

	var form NotificationBatchForm

	var page int
	vars := mux.Vars(r)
	pageStr := vars["page"]
	if pageStr != "" {
		page, err = strconv.Atoi(vars["page"])
		if err != nil {
			api.SendErrorJSON(w, r, http.StatusBadRequest, err, "validation error", api.ErrValidationData)
			return
		}
	} else {
		page = 1
	}

	form.Page = page

	_, err = form.Validate()
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "validation error", api.ErrValidationData)
		return
	}

	limit := 10

	var items []models.Notification

	response := NotificationBatch{Total: 0, Items: &items}

	total, err := s.readRepository.GetTotalForUser(currentUser.ID)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get total notifications for user", api.ErrInternal)
		return
	}

	if total == 0 {
		api.SendDataJSON(w, r, http.StatusOK, response)
		return
	}

	response.Total = total

	items, err = s.readRepository.GetBatchForUser(currentUser.ID, page, limit)
	if err != nil {
		if _, ok := err.(repositories.ErrNotFound); !ok {
			api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get notifications for user", api.ErrInternal)
			return
		}

		api.SendDataJSON(w, r, http.StatusOK, response)
		return
	}

	api.SendDataJSON(w, r, http.StatusOK, response)
}

func (s *NotificationController) notificationPush(w http.ResponseWriter, r *http.Request) {
	if s.PushChannel == nil {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, nil, "PushChannel is nil", api.ErrInternal)
		return
	}

	currentUser, err := GetAuthorizedUser(s.authService, r)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusUnauthorized, err, "can't auth user by token", api.ErrUserNotFound)
		return
	}

	var form NotificationPushForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse", api.ErrDecode)
		return
	}

	_, err = form.Validate()
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "validation error", api.ErrValidationData)
		return
	}

	notification, err := s.repository.GetById(form.ID)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get notification", api.ErrInternal)
		return
	}

	pushMessage := services.PushNotificationMessage{Notification: notification, TargetUserIds: []string{currentUser.ID}}
	s.PushChannel <- pushMessage

	api.SendDataJSON(w, r, http.StatusOK, nil)
}
