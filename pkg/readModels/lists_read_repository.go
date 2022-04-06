package readModels

import (
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"shopingList/pkg/models"
	"shopingList/pkg/repositories"
	"strings"
)

const TableLists = "sl_item_list"

type ListsReadRepository struct {
	DB *sql.DB
}

// Вернуть список по ID и владельцу
func (s *ListsReadRepository) GetListForIdAndOwner(id string, ownerId string) (models.List, error) {
	db := s.DB
	row := db.QueryRow(
		`SELECT 
       id, 
       owner_id, 
       name, 
       is_template,
       UNIX_TIMESTAMP(created_at), 
       UNIX_TIMESTAMP(updated_at), 
       UNIX_TIMESTAMP(received_at), 
       is_deleted
		FROM sl_item_list
		WHERE id =? AND owner_id =? `,
		id, ownerId)

	var l models.List
	err := row.Scan(&l.ID, &l.OwnerID, &l.Name, &l.IsTemplate, &l.CreatedAt, &l.UpdatedAt, &l.ReceivedAt, &l.IsDeleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.List{}, repositories.ErrNotFound{}
		}

		return models.List{}, err
	}
	return l, nil
}

// Вернуть обновленные списки пользователя
func (s *ListsReadRepository) GetUpdatedListsForOwner(ownerID string, receivedAt int64) ([]models.List, error) {
	db := s.DB
	rows, err := db.Query(
		s.getSelectPartSql()+` WHERE l.owner_id=? AND l.received_at >= FROM_UNIXTIME(?)`,
		ownerID, receivedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repositories.ErrNotFound{}
		}

		return nil, err
	}
	defer rows.Close() // nolint errcheck

	return s.listRowsToArray(rows)
}

func (s *ListsReadRepository) getSelectPartSql() string {
	return `SELECT 
			l.id, 
			l.owner_id, 
			l.name,
			l.is_template,
			UNIX_TIMESTAMP(l.created_at), 
			UNIX_TIMESTAMP(l.updated_at), 
			UNIX_TIMESTAMP(l.received_at),
			l.is_deleted 
			FROM sl_item_list AS l `
}

// Вернуть списки, принадлежащие пользователю по ID списков и ID владельца
// Возврат по ID и Owner нужен потому что уникальность списка должна быть в паре ID + ownerId
func (s *ListsReadRepository) GetListsForIdsAndOwner(listIds []string, ownerId string) ([]models.List, error) {
	db := s.DB

	if listIds == nil {
		return nil, sql.ErrNoRows
	}

	var args []interface{}
	for _, id := range listIds {
		args = append(args, id)
	}

	sqlQuery := s.getSelectPartSql() + `
			WHERE l.id IN (?` + strings.Repeat(`,?`, len(args)-1) + `) AND l.owner_id =?`

	stmt, err := db.Prepare(sqlQuery)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repositories.ErrNotFound{}
		}

		return nil, err
	}

	args = append(args, ownerId)
	rows, err := stmt.Query(args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.listRowsToArray(rows)
}

// Вернуть обновленные списки, пошаренные на пользователя
// Учитывается дата шаринга, чтобы возвращались пошаренные списки
// с датой обновления более старой, чем дата шаринга.
func (s *ListsReadRepository) GetUpdatedListsSharedToUser(toUserID string, receivedAt int64) ([]models.List, error) {
	db := s.DB
	rows, err := db.Query(
		s.getSelectPartSql()+`LEFT JOIN sl_shared_lists AS s ON (l.id = s.list_id)
			WHERE to_user_id=? AND (s.received_at >= FROM_UNIXTIME(?) OR l.received_at >= FROM_UNIXTIME(?)) `,
		toUserID, receivedAt, receivedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repositories.ErrNotFound{}
		}

		return nil, err
	}
	defer rows.Close()

	return s.listRowsToArray(rows)
}

// Вернуть списки, пошаренные на пользователя по ID списков
func (s *ListsReadRepository) GetListsSharedForUserForIds(listIds []string, userId string) ([]models.List, error) {
	db := s.DB

	if listIds == nil {
		return nil, sql.ErrNoRows
	}

	var args []interface{}
	for _, id := range listIds {
		args = append(args, id)
	}

	sqlQuery := s.getSelectPartSql() + `LEFT JOIN sl_shared_lists AS s ON (l.id = s.list_id)
			WHERE l.id IN (?` + strings.Repeat(`,?`, len(args)-1) + `) AND s.to_user_id =?`

	stmt, err := db.Prepare(sqlQuery)

	if err != nil {
		return nil, err
	}

	args = append(args, userId)
	rows, err := stmt.Query(args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repositories.ErrNotFound{}
		}

		return nil, err
	}
	defer rows.Close()

	return s.listRowsToArray(rows)
}

func (s *ListsReadRepository) GetActiveListsForUser(userId string) ([]models.List, error) {
	var lists []models.List

	selDataset := s.getSelectDataset()
	err := selDataset.Where(goqu.Ex{"owner_id": userId}, goqu.Ex{"is_deleted": false}).
		Order(goqu.I("updated_at").Asc()).
		ScanStructs(&lists)

	if err != nil {
		return nil, errors.Wrap(err, "Error get active lists by user_id")
	}

	return lists, nil
}

func (r *ListsReadRepository) getGoquDB() *goqu.Database {
	return goqu.New("mysql", r.DB)
}

func (r *ListsReadRepository) getSelectDataset() *goqu.SelectDataset {
	db := r.getGoquDB()

	return db.From(TableLists).Select("id", "owner_id", "name", "is_template",
		goqu.L("UNIX_TIMESTAMP(`created_at`)").As("created_at"),
		goqu.L("UNIX_TIMESTAMP(`updated_at`)").As("updated_at"),
		goqu.L("UNIX_TIMESTAMP(`updated_at`)").As("received_at"),
		"is_deleted")
}

func (s *ListsReadRepository) listRowsToArray(rows *sql.Rows) ([]models.List, error) {
	var lists []models.List

	for rows.Next() {
		var l models.List
		err := rows.Scan(&l.ID, &l.OwnerID, &l.Name, &l.IsTemplate, &l.CreatedAt, &l.UpdatedAt, &l.ReceivedAt, &l.IsDeleted)
		if err != nil {
			return nil, err
		}
		lists = append(lists, l)
	}
	return lists, nil
}
