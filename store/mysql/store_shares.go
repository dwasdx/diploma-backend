package mysql

import (
	"database/sql"
	"shopingList/pkg/models"
	"strings"
)

// Вернуть список шарингов для владельца
func (s *DataStore) GetUpdatedSharesForOwner(ownerID string, updatedAt int64) ([]models.ListShare, error) {
	db := s.db
	rows, err := db.Query(
		`SELECT s.id,
				s.list_id, 
       			s.to_user_id, 
       			s.owner_id,
       			s.status, 
       			UNIX_TIMESTAMP(s.created_at), 
       			UNIX_TIMESTAMP(s.updated_at), 
       			s.is_deleted
			FROM sl_shared_lists AS s
			LEFT JOIN sl_item_list AS l ON (s.list_id = l.id)
			WHERE l.owner_id=? AND UNIX_TIMESTAMP(s.updated_at) >=?`,
		ownerID, updatedAt,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close() // nolint errcheck

	return shareRowsToArray(rows)
}

// Вернуть список шарингов предназначенных для получателя
func (s *DataStore) GetUpdatedSharesToUser(toUserID string, updatedAt int64) ([]models.ListShare, error) {
	db := s.db
	rows, err := db.Query(
		`SELECT s.id,
				s.list_id, 
       			s.to_user_id, 
       			s.owner_id,
       			s.status, 
       			UNIX_TIMESTAMP(s.created_at), 
       			UNIX_TIMESTAMP(s.updated_at), 
       			s.is_deleted
			FROM sl_shared_lists AS s
			LEFT JOIN sl_item_list AS l ON (s.list_id = l.id)
			WHERE to_user_id=? AND UNIX_TIMESTAMP(s.updated_at) >=?`,
		toUserID, updatedAt,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close() // nolint errcheck

	return shareRowsToArray(rows)
}

func shareRowsToArray(rows *sql.Rows) ([]models.ListShare, error) {
	var shares []models.ListShare

	for rows.Next() {
		var share models.ListShare
		err := rows.Scan(
			&share.ID,
			&share.ListID,
			&share.ToUserID,
			&share.OwnerID,
			&share.Status,
			&share.CreatedAt,
			&share.UpdatedAt,
			&share.IsDeleted,
		)

		if err != nil {
			return nil, err
		}
		shares = append(shares, share)
	}

	return shares, nil
}

// Вернуть шаринги, пошаренные на пользователя по ID списков
func (s *DataStore) GetSharesForUserForListIds(listIds []string, userId string) ([]models.ListShare, error) {
	db := s.db

	if listIds == nil {
		return nil, sql.ErrNoRows
	}

	var args []interface{}
	for _, id := range listIds {
		args = append(args, id)
	}

	sqlQuery := `SELECT s.id,
				s.list_id, 
       			s.to_user_id, 
       			s.status, 
       			UNIX_TIMESTAMP(s.created_at), 
       			UNIX_TIMESTAMP(s.updated_at), 
       			s.is_deleted
			FROM sl_shared_lists AS s
			LEFT JOIN sl_item_list AS l ON (s.list_id = l.id)
			WHERE s.list_id IN (?` + strings.Repeat(`,?`, len(args)-1) + `) AND s.to_user_id=?`

	stmt, err := db.Prepare(sqlQuery)

	if err != nil {
		return nil, err
	}

	args = append(args, userId)
	rows, err := stmt.Query(args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return shareRowsToArray(rows)
}
