package repositories

import (
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/doug-martin/goqu/v9/exec"
	"github.com/pkg/errors"
	"shopingList/pkg"
	"shopingList/pkg/tg"
	"time"
)

const TableTgUsers = "sl_tg_users"

type TgUsersRepository struct {
	db *sql.DB
}

func NewTgUsersRepository(db *sql.DB) *TgUsersRepository {
	return &TgUsersRepository{db: db}
}

func (r *TgUsersRepository) Get(id int) (*tg.TgUser, error) {
	tgUser := tg.TgUser{}

	_, err := r.getSelectDataset().Where(goqu.Ex{"tg_id": id}).ScanStruct(&tgUser)
	if err != nil {
		return nil, errors.Wrap(err, "Error get tguser by id")
	}

	if tgUser.TgID == 0 {
		return nil, pkg.ErrNotFoundInStorage
	}

	return &tgUser, nil
}

func (r *TgUsersRepository) Save(user *tg.TgUser) error {
	db := r.getGoquDB()

	count, err := db.From(TableTgUsers).Where(goqu.Ex{"tg_id": user.TgID}).Count()
	if err != nil {
		return errors.Wrap(err, "Error check exist")
	}

	user.UpdatedAt = time.Now().UTC().Unix()

	record := goqu.Record{"tg_id": user.TgID, "username": user.Username, "user_id": user.UserId}
	updated_at := time.Now().UTC().Format("2006-01-02 15:04:05")
	record["updated_at"] = updated_at

	var executor exec.QueryExecutor

	if count == 0 {
		record["created_at"] = updated_at
		executor = db.Insert(TableTgUsers).Rows(record).Executor()
	} else {
		executor = db.Update(TableTgUsers).Where(goqu.Ex{"tg_id": user.TgID}).Set(record).Executor()
	}

	if _, err = executor.Exec(); err != nil {
		return errors.Wrap(err, "Error save tgUser")
	}

	return nil
}

func (r *TgUsersRepository) getGoquDB() *goqu.Database {
	return goqu.New("mysql", r.db)
}

func (r *TgUsersRepository) getSelectDataset() *goqu.SelectDataset {
	db := r.getGoquDB()

	return db.From(TableTgUsers).Select("tg_id", "username", "user_id",
		goqu.L("UNIX_TIMESTAMP(`created_at`)").As("created_at"),
		goqu.L("UNIX_TIMESTAMP(`updated_at`)").As("updated_at"))
}
