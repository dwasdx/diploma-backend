package repositories

import (
	"database/sql"
	"errors"
	"shopingList/pkg/models"
)

type NotificationsRepository struct {
	db models.DB
}

func NewNotificationsRepository(db models.DB) NotificationsRepository {
	if db == nil {
		panic("db param is nil")
	}

	return NotificationsRepository{db: db}
}

func (s *NotificationsRepository) GetById(id string) (models.Notification, error) {
	row := s.db.QueryRow(
		`SELECT id,
       			type, 
       			message,
       			user_id, 
       			user_phone,
       			list_id, 
       			item_id, 
				target_user_id,
       			UNIX_TIMESTAMP(created_at) 
		FROM `+models.NotificationTableName+`
		WHERE id = ?`, id)

	var model models.Notification

	err := row.Scan(
		&model.ID, &model.TypeNotification, &model.Message, &model.UserId, &model.UserPhone,
		&model.ListId, &model.ItemId, &model.TargetUserId, &model.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Notification{}, ErrNotFound{}
		}

		return models.Notification{}, err
	}

	return model, nil
}

func (s *NotificationsRepository) Create(form *models.NotificationCreateForm) error {
	_, err := s.db.Exec(`INSERT INTO `+models.NotificationTableName+` (
                    id, type, message, user_id, user_phone, list_id, item_id, target_user_id
                    ) 
		VALUES (UUID(), ?, ?, ?, ?, ?, ?, ?)`,
		form.TypeNotification, form.Message, form.UserId, form.UserPhone, form.ListId, form.ItemId.SqlValue(), form.TargetUserId)

	if err != nil {
		return errors.New("Error insert notification; " + err.Error())
	}

	return nil
}
