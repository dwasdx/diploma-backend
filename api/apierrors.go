package api

// API errors
const (
	ErrInternal           = 1 // any internal error
	ErrDecode             = 2
	ErrNoPermission       = 3
	ErrCreateUser         = 4
	ErrLogin              = 5
	ErrUserAlreadyExists  = 6
	ErrUserNotFound       = 7
	ErrUserHasDeleted     = 8
	ErrUserIsNotActivated = 9
	ErrInvalidCode        = 10
	ErrCredentials        = 11
	ErrInvalidToken       = 12
	ErrValidationData     = 13
	ErrLimitLogin         = 14 // Ограничение на количество авторизации в период времени
)
