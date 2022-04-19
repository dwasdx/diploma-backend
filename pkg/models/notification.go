package models

import (
	"github.com/asaskevich/govalidator"
	"strconv"
)

const NotificationTableName = string("sl_notifications")

const (
	NotificationTypeListInvite      = 1  // Приглашение в список
	NotificationTypeListJoining     = 2  // Присоединение к списку
	NotificationTypeListDetachment  = 3  // Отсоединение от списка
	NotificationTypeGoodsCreate     = 4  // Добавить новый товар
	NotificationTypeGoodsCheck      = 5  // Отметка товара
	NotificationTypeGoodsUncheck    = 6  // Снять отметку с товара
	NotificationTypeGoodsChange     = 7  // Изменение товара
	NotificationTypeGoodsDelete     = 8  // Удаление товара
	NotificationTypeListShareDelete = 9  // Удаление шаринга
	NotificationTypeListDelete      = 10 // Удаление списка
)

type NotificationType int

type Notification struct {
	TypeNotification NotificationType `json:"type"  valid:"uuid,required"`
	ID               string           `json:"id" valid:"uuid,required"`
	Message          string           `json:"message" valid:"stringlength(1|255),required"`
	UserId           string           `json:"user_id" valid:"uuid"`
	UserPhone        int64            `json:"user_phone" valid:"stringlength(0|10)"`
	ListId           string           `json:"list_id" valid:"uuid"`
	ItemId           NullString       `json:"item_id" valid:"uuid"`
	TargetUserId     string           `json:"-"`
	CreatedAt        int64            `json:"created_at" valid:"int,required"`
}

type NotificationCreateForm struct {
	TypeNotification NotificationType `json:"type"  valid:"uuid,required"`
	Message          string           `json:"name" valid:"stringlength(1|255),required"`
	UserId           string           `json:"user_id" valid:"uuid"`
	UserPhone        int64            `json:"user_phone" valid:"stringlength(0|10)"`
	ListId           string           `json:"list_id" valid:"uuid"`
	ItemId           NullString       `json:"item_id" valid:"uuid"`
	TargetUserId     string           `json:"-" valid:"uuid,required"`
}

func (s *NotificationCreateForm) Validate() (bool, error) {
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s NotificationCreateForm) GetTypeNotification() string {
	return strconv.Itoa(int(s.TypeNotification))
}

func (s NotificationCreateForm) GetUserPhone() string {
	return strconv.FormatInt(s.UserPhone, 10)
}

func (s NotificationCreateForm) GetUserId() string {
	return s.UserId
}

func (s NotificationCreateForm) GetListId() string {
	return s.ListId
}

func (s NotificationCreateForm) GetItemId() string {
	return s.ItemId.String
}

func (s NotificationCreateForm) GetMessage() string {
	return s.Message
}

func (s Notification) GetTypeNotification() string {
	return strconv.Itoa(int(s.TypeNotification))
}

func (s Notification) GetUserPhone() string {
	return strconv.FormatInt(s.UserPhone, 10)
}

func (s Notification) GetUserId() string {
	return s.UserId
}

func (s Notification) GetListId() string {
	return s.ListId
}

func (s Notification) GetItemId() string {
	return s.ItemId.String
}

func (s Notification) GetMessage() string {
	return s.Message
}
