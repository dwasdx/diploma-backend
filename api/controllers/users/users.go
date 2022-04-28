package users

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
	"shopingList/api"
	auth2 "shopingList/api/auth"
	"shopingList/pkg"
	"shopingList/pkg/forms"
	"shopingList/pkg/models"
	"shopingList/pkg/repositories"
	"shopingList/pkg/services/sms"
	"shopingList/store"
	"strconv"
	"time"
)

type UserController struct {
	authService   *auth2.Service
	dataService   store.DataService
	smsService    sms.SmsService
	codeGenerator models.UserAuthCodeGenerator
	limiter       LoginLimiter
	debugPhones   map[int64]bool
}

func NewUserController(
	authService *auth2.Service,
	dataService store.DataService,
	smsService sms.SmsService,
	limiter LoginLimiter,
	codeGenerator models.UserAuthCodeGenerator) *UserController {
	return &UserController{
		authService:   authService,
		dataService:   dataService,
		smsService:    smsService,
		limiter:       limiter,
		codeGenerator: codeGenerator}
}

// Routes returns slice of server routes
func (s *UserController) Routes() []api.Route {
	return []api.Route{
		{
			Name:   "CreateUser",
			Method: "POST",
			Path:   "/user/create",
			Func:   s.createUser,
		},
		{
			Name:   "PhoneConfirm",
			Method: "POST",
			Path:   "/auth/phone/confirm",
			Func:   s.phoneConfirm,
		},
		{
			Name:   "UserLogin",
			Method: "POST",
			Path:   "/auth/phone",
			Func:   s.login,
		},
		{
			Name:   "NextTimeLogin",
			Method: "GET",
			Path:   "/auth/phone/{phone}/checkLimit",
			Func:   s.checkAllowLogin,
		},
	}
}

func (s *UserController) createUser(w http.ResponseWriter, r *http.Request) {
	var form forms.RegForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse", api.ErrDecode)
		return
	}

	_, err := form.Validate()
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "validation error", api.ErrValidationData)
		return
	}

	usersReadRepository := s.dataService.GetUsersReadRepository()
	usersRepository := s.dataService.GetUsersRepository(nil)

	user, err := usersReadRepository.GetUserByPhone(form.Phone)
	if err != nil && err != pkg.ErrNotFoundInStorage {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get user by phone", api.ErrInternal)
		return
	}

	if user != nil && user.IsActivated {
		err := errors.New("user is exist")
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't create user. User is already exists", api.ErrUserAlreadyExists)
		return
	}

	if user == nil {
		user, err = usersRepository.CreateUser(form)
		if err != nil {
			api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't create user", api.ErrCreateUser)
			return
		}
	} else {
		user.Name = form.Username
		user.Email.String = form.Email
		user.UpdatedAt = time.Now().UTC().Unix()

		err = usersRepository.UpdateUser(user)
		if err != nil {
			api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't update user", api.ErrCreateUser)
			return
		}
	}

	if user == nil {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't create user, user = nil", api.ErrCreateUser)
		return
	}

	err = s.sendAuthCodeToUser(user, usersRepository)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "Error send sms with auth code", api.ErrInternal)
		return
	}

	api.SendDataJSON(w, r, http.StatusOK, user)
}

func (s *UserController) login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't login", api.ErrLogin)
		return
	}

	_, err := form.Validate()
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "validation error", api.ErrValidationData)
		return
	}

	usersRepository := s.dataService.GetUsersRepository(nil)
	usersReadRepository := s.dataService.GetUsersReadRepository()

	user, err := usersReadRepository.GetUserByPhone(form.Phone)
	if err != nil {
		if err == pkg.ErrNotFoundInStorage {
			api.SendErrorJSON(w, r, http.StatusBadRequest, errors.New("user doesn`t found by phone"), "can't get user", api.ErrUserNotFound)
			return
		}

		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't get user", api.ErrUserNotFound)
		return
	}

	isAllow := true
	isDebugPhone := s.isDebugPhone(form.Phone)

	if !isDebugPhone {
		isAllow, err = s.limiter.IsAllowLogin(form.Phone)
		if err != nil {
			api.SendErrorJSON(w, r, http.StatusInternalServerError, err,
				"Error in check limit", api.ErrInternal)
			return
		}
	}

	if !isAllow {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err,
			"Authorization temporarily prohibited due to time limit", api.ErrLimitLogin)
		return
	}

	nextTime, err := s.limiter.RememberOperation(form.Phone)

	if !isDebugPhone {
		err = s.sendAuthCodeToUser(user, usersRepository)
		if err != nil {
			api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "Error send sms with auth code", api.ErrInternal)
			return
		}
	}

	api.SendDataJSON(w, r, http.StatusOK, api.JSON{"user_id": user.ID, "next_time_request": nextTime})
}

func (s *UserController) isDebugPhone(phone int64) bool {
	_, isDebugPhone := s.debugPhones[phone]
	return isDebugPhone
}

func (s *UserController) checkAllowLogin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	phoneStr := vars["phone"]

	phone, err := strconv.ParseInt(phoneStr, 10, 64)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "Error parse phone", api.ErrInternal)
		return
	}

	result, err := s.limiter.NextTimeRequest(phone)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "Error get next login", api.ErrInternal)
		return
	}

	api.SendDataJSON(w, r, http.StatusOK, result)
}

func (s *UserController) phoneConfirm(w http.ResponseWriter, r *http.Request) {
	var form PhoneConfirm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse", api.ErrDecode)
		return
	}

	_, err := form.Validate()
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "validation error", api.ErrValidationData)
		return
	}

	usersRepository := s.dataService.GetUsersRepository(nil)

	user, err := usersRepository.GetUser(form.UserID)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't find user", api.ErrUserNotFound)
		return
	}

	isDebugPhone := s.isDebugPhone(user.Phone)

	// TODO Добавить проверку на время жизни токена
	if !isDebugPhone {
		if user.Code.String != form.Code {
			api.SendErrorJSON(w, r, http.StatusBadRequest, errors.New("invalid code"), "can't verify code", api.ErrInvalidCode)
			return
		}
	}

	user.ClearAuthCode()

	if !user.IsActivated {
		user.IsActivated = true
	}

	err = usersRepository.UpdateUser(user)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't update user", api.ErrInternal)
		return
	}

	claims := auth2.Claims{
		UserUUID: user.ID,
		StandardClaims: &jwt.StandardClaims{
			Subject:   "access",
			IssuedAt:  time.Now().UTC().Unix(),
			ExpiresAt: time.Now().Add(24 * time.Hour * 365).UTC().Unix(),
			Issuer:    "shopping_list_backend",
		},
	}

	tokenString, err := s.authService.CreateJWT(claims)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't create token", api.ErrInternal)
		return
	}

	resp := api.JSON{"user": user, "user_id": user.ID, "access_token": tokenString, "expires_at": claims.ExpiresAt}

	api.SendDataJSON(w, r, http.StatusOK, resp)
}

func (s *UserController) sendAuthCodeToUser(user *models.User, repository repositories.UsersRepository) error {
	user.RegenerateAuthCode(s.codeGenerator)
	if user.Code.String == "" {
		return errors.New("empty activation code")
	}

	err := repository.UpdateUser(user)
	if err != nil {
		return errors.New("can't update user; " + err.Error())
	}

	message := "Shoppinglist. Код подтверждения: " + user.Code.String

	if s.smsService != nil {
		err := s.smsService.Send(user.Phone, message)
		if err != nil {
			return errors.New("can't send sms; " + err.Error())
		}
	}

	return nil
}

func (s *UserController) SetDebugPhones(phones []int64) {
	s.debugPhones = make(map[int64]bool, 0)

	for _, phone := range phones {
		s.debugPhones[phone] = true
	}
}
