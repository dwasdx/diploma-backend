package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"math"
	"net/http"
	"shopingList/api"
	"shopingList/pkg"
	"shopingList/pkg/forms"
	"strconv"

	"shopingList/api/auth"
	"shopingList/pkg/models"
	"shopingList/store"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/mux"
)

// Private API-service
type Private struct {
	dataService store.DataService
}

// NewPrivate returns new instance of private API
func NewPrivate(dataService store.DataService) *Private {
	return &Private{dataService: dataService}
}

func GetAuthorizedUser(authService *auth.Service, r *http.Request) (*models.User, error) {
	currentUser, err := authService.UserInfo(r)
	if err != nil {
		return nil, err
	}

	return currentUser, nil
}

// Routes returns slice of server routes
func (s *Private) Routes() []api.Route {
	return []api.Route{
		{
			Name:   "User",
			Method: "GET",
			Path:   "/user/by-phone/{phone}",
			Func:   s.getOrCreateUserByPhone,
		},
		{
			Name:   "UserByPhones",
			Method: "POST",
			Path:   "/users/check-registation",
			Func:   s.getUsersActivity,
		},
	}
}

func (s *Private) getOrCreateUserByPhone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	phoneString := vars["phone"]

	err := validation.Validate(phoneString, validation.Required, is.Digit, validation.Length(11, 15))
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "wrong format number phone", api.ErrDecode)
		return
	}

	phone, err := strconv.ParseInt(phoneString, 10, 64)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "error convert phone number", api.ErrDecode)
		return
	}

	readRepository := s.dataService.GetUsersReadRepository()
	repository := s.dataService.GetUsersRepository(nil)

	user, err := readRepository.GetUserByPhone(phone)
	if err != nil {
		if err == pkg.ErrNotFoundInStorage {
			userCreated, err2 := repository.CreateUser(forms.RegForm{Phone: phone})

			if err2 != nil {
				api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't create user", api.ErrInternal)
				return
			}

			user = userCreated
		} else {
			api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get user", api.ErrInternal)
			return
		}
	}

	api.SendDataJSON(w, r, http.StatusOK, models.CreateExternalUser(user))
	return
}

const maxPhones = 1000
const chunkSize = 200

type CheckUsersActivity struct {
	Phones []int64 `json:"phones" valid:"required"`
}

func (s *CheckUsersActivity) Validate() (bool, error) {
	if _, err := govalidator.ValidateStruct(s); err != nil {
		return false, err
	}

	if len(s.Phones) > maxPhones {
		return false, errors.New(fmt.Sprintf("phone numbers more than %d", maxPhones))
	}

	return true, nil
}

func (s *Private) getUsersActivity(w http.ResponseWriter, r *http.Request) {
	form := CheckUsersActivity{}

	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't validate request", api.ErrValidationData)
		return
	}

	if _, err := form.Validate(); err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't validate form request", api.ErrValidationData)
		return
	}

	phonesChunked := make([][]int64, 0)

	if len(form.Phones) > chunkSize {
		var first, last int

		for i := 0; i < int(math.Ceil(float64(len(form.Phones)+1)/chunkSize)); i++ {
			first = i * chunkSize
			last = i*chunkSize + chunkSize

			if last >= len(form.Phones) {
				last = len(form.Phones) - 1
			}

			if first >= last {
				break
			}

			phonesChunked = append(phonesChunked, form.Phones[first:last])
		}
	} else {
		phonesChunked = append(phonesChunked, form.Phones)
	}

	responseData := make(map[int64]bool)
	for _, phone := range form.Phones {
		responseData[phone] = false
	}

	for _, phoneChunk := range phonesChunked {
		if err := s.loadAndCheckPhones(phoneChunk, responseData); err != nil {
			api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "error search users by phone array", api.ErrInternal)
			return
		}
	}

	api.SendDataJSON(w, r, http.StatusOK, responseData)
	return
}

func (s *Private) loadAndCheckPhones(phones []int64, responseData map[int64]bool) error {
	usersRepositories := s.dataService.GetUsersReadRepository()
	users, err := usersRepositories.GetUsersByPhones(phones...)
	if err != nil {
		if err == pkg.ErrNotFoundInStorage {
			return nil
		}

		return err
	}

	for _, user := range users {
		if user.IsActivated {
			responseData[user.Phone] = true
		}
	}

	return nil
}
