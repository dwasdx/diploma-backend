package readModels

import (
	"database/sql"
	"fmt"
	"shopingList/pkg/models"
	"shopingList/pkg/repositories"
)

type NotificationsReadRepository struct {
	db *sql.DB
}

func NewNotificationsReadRepository(db *sql.DB) NotificationsReadRepository {
	if db == nil {
		panic("DB is nil")
	}

	return NotificationsReadRepository{db: db}
}

func (s *NotificationsReadRepository) GetTotalForUser(userId string) (int, error) {
	rows, err := s.db.Query(
		`SELECT count(*) 
			FROM `+models.NotificationTableName+`
			WHERE target_user_id =?`,
		userId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, repositories.ErrNotFound{}
		}

		return 0, err
	}
	defer rows.Close()

	count := 0

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (s *NotificationsReadRepository) GetBatchForUser(userId string, page int, limit int) ([]models.Notification, error) {
	var items []models.Notification

	if page <= 0 {
		page = 1
	}

	start := page*limit - limit

	sqlLimit := fmt.Sprintf("LIMIT %d, %d", start, limit)

	rows, err := s.db.Query(
		`SELECT id, type, message, user_id, user_phone, list_id, item_id, UNIX_TIMESTAMP(created_at) 
			FROM `+models.NotificationTableName+`
			WHERE target_user_id =? ORDER BY created_at DESC `+sqlLimit,
		userId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repositories.ErrNotFound{}
		}

		return nil, err
	}
	defer rows.Close()

	items, err = s.scanRows(rows)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *NotificationsReadRepository) scanRows(rows *sql.Rows) ([]models.Notification, error) {
	var items []models.Notification

	for rows.Next() {
		var model models.Notification
		err := rows.Scan(&model.ID, &model.TypeNotification, &model.Message, &model.UserId, &model.UserPhone, &model.ListId, &model.ItemId, &model.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, model)
	}
	return items, nil
}
