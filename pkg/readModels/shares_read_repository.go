package readModels

import (
	"database/sql"
	"shopingList/pkg"
	"shopingList/pkg/models"
	"strings"
)

const SharesTableName = "sl_shared_lists"

type SharesReadRepository struct {
	db *sql.DB
}

func NewSharesReadRepository(db *sql.DB) SharesReadRepository {
	if db == nil {
		panic("DB is nil")
	}

	return SharesReadRepository{db: db}
}

// Вернуть объект шаринга
func (s *SharesReadRepository) GetShare(id string, ownerId string) (*models.ListShare, error) {
	row := s.db.QueryRow(
		`SELECT s.id,
				s.list_id, 
       			s.to_user_id,
       			s.owner_id,
       			s.status, 
       			UNIX_TIMESTAMP(s.created_at), 
       			UNIX_TIMESTAMP(s.updated_at), 
				UNIX_TIMESTAMP(s.received_at),
       			s.is_deleted
			FROM `+SharesTableName+` AS s
			WHERE s.id =? AND s.owner_id=?`,
		id, ownerId)

	var share models.ListShare
	err := row.Scan(&share.ID, &share.ListID, &share.ToUserID, &share.OwnerID, &share.Status,
		&share.CreatedAt, &share.UpdatedAt, &share.ReceivedAt, &share.IsDeleted)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.ErrNotFoundInStorage
		}

		return nil, err
	}

	return &share, nil
}

// Вернуть объект шаринга
func (s *SharesReadRepository) GetShareForUser(id string, toUserId string) (*models.ListShare, error) {
	row := s.db.QueryRow(
		`SELECT id,
				list_id, 
       			to_user_id,
       			owner_id,
       			status, 
       			UNIX_TIMESTAMP(created_at), 
       			UNIX_TIMESTAMP(updated_at), 
				UNIX_TIMESTAMP(received_at),
       			is_deleted
			FROM `+SharesTableName+`
			WHERE id =? AND to_user_id=?`,
		id, toUserId)

	var share models.ListShare
	err := row.Scan(&share.ID, &share.ListID, &share.ToUserID, &share.OwnerID, &share.Status,
		&share.CreatedAt, &share.UpdatedAt, &share.ReceivedAt, &share.IsDeleted)

	if err != nil {
		return nil, err
	}

	return &share, nil
}

// Вернуть список шарингов для владельца
func (s *SharesReadRepository) GetUpdatedSharesForOwner(ownerID string, receivedAt int64) ([]models.ListShare, error) {
	rows, err := s.db.Query(
		`SELECT s.id,
				s.list_id, 
       			s.to_user_id, 
       			s.owner_id,
       			s.status, 
       			UNIX_TIMESTAMP(s.created_at), 
       			UNIX_TIMESTAMP(s.updated_at), 
				UNIX_TIMESTAMP(s.received_at),
       			s.is_deleted
			FROM `+SharesTableName+`  AS s
			LEFT JOIN sl_item_list AS l ON (s.list_id = l.id)
			WHERE l.owner_id=? AND s.received_at >= FROM_UNIXTIME(?)`,
		ownerID, receivedAt,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close() // nolint errcheck

	return shareRowsToArray(rows)
}

// Вернуть список шарингов предназначенных для получателя
func (s *SharesReadRepository) GetUpdatedSharesToUser(toUserID string, receivedAt int64) ([]models.ListShare, error) {
	db := s.db
	rows, err := db.Query(
		`SELECT s.id,
				s.list_id, 
       			s.to_user_id, 
       			s.owner_id,
       			s.status, 
       			UNIX_TIMESTAMP(s.created_at), 
       			UNIX_TIMESTAMP(s.updated_at), 
				UNIX_TIMESTAMP(s.received_at),
       			s.is_deleted
			FROM `+SharesTableName+`  AS s
			LEFT JOIN sl_item_list AS l ON (s.list_id = l.id)
			WHERE to_user_id=? AND s.received_at >= FROM_UNIXTIME(?)`,
		toUserID, receivedAt,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close() // nolint errcheck

	return shareRowsToArray(rows)
}

// Вернуть шаринги, пошаренные на пользователя по ID списков
func (s *SharesReadRepository) GetSharesForUserForListIds(listIds []string, userId string) ([]models.ListShare, error) {
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
				s.owner_id,
       			s.status, 
       			UNIX_TIMESTAMP(s.created_at), 
       			UNIX_TIMESTAMP(s.updated_at), 
				UNIX_TIMESTAMP(s.received_at),
       			s.is_deleted
			FROM ` + SharesTableName + `  AS s
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

// Вернуть массив ID акцентованных пользователей списка
func (s *SharesReadRepository) GetAcceptedUserIdsFromSharedList(listId string, ownerId string) ([]string, error) {
	var userIds []string

	rows, err := s.db.Query(
		`SELECT s.to_user_id FROM `+SharesTableName+`  AS s
			WHERE list_id=? AND owner_id=? AND status=? AND is_deleted = 0`,
		listId, ownerId, models.ShareStatusAccepted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	defer rows.Close() // nolint errcheck

	for rows.Next() {
		var id string
		err := rows.Scan(
			&id,
		)

		if err != nil {
			return nil, err
		}
		userIds = append(userIds, id)
	}

	return userIds, nil
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
			&share.ReceivedAt,
			&share.IsDeleted,
		)

		if err != nil {
			return nil, err
		}
		shares = append(shares, share)
	}

	return shares, nil
}
