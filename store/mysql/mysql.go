package mysql

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"log"
	"shopingList/pkg/readModels"
	"shopingList/pkg/repositories"
)

// DataStore struct for working with mysql
type DataStore struct {
	db *sql.DB
}

// NewDataStore returns new MySQL store service with connection pool
func NewDataStore(db *sql.DB) *DataStore {
	return &DataStore{db: db}
}

func prepareDatabase(db *sql.DB, name string) error {
	log.Println("[INFO]: prepare database")
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec("USE " + name)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
        CREATE TABLE IF NOT EXISTS  sl_users (
            id varchar(36) NOT NULL,
            name varchar(100) NOT NULL DEFAULT '',
            email varchar(100) NOT NULL DEFAULT '',
            phone bigint(11) NOT NULL,
            code varchar(8) NOT NULL DEFAULT '',
            created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            is_activated tinyint(1) NOT NULL DEFAULT '0',
            is_deleted tinyint(1) NOT NULL DEFAULT '0',
        PRIMARY KEY (id),
        UNIQUE KEY sl_users_phone_unique (phone),
        KEY phone (phone)
         ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS sl_item_list (
  			id varchar(36) NOT NULL,
            owner_id varchar(36) NOT NULL,
            name varchar(100) NOT NULL DEFAULT '',
            created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            is_deleted tinyint(1) NOT NULL DEFAULT '0',
        PRIMARY KEY (id),
        KEY fk_owner_user (owner_id),
        CONSTRAINT fk_owner_user FOREIGN KEY (owner_id) 
            REFERENCES shoppinglist.sl_users (id) ON DELETE CASCADE ON UPDATE NO ACTION
		) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
	`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS sl_item (
  			id varchar(36) NOT NULL,
			name varchar(140) NOT NULL DEFAULT '',
            value varchar(50) NOT NULL DEFAULT '',
            is_marked tinyint(1) NOT NULL DEFAULT '0',
            user_marked_id varchar(36) DEFAULT NULL,
            list_id varchar(36) NOT NULL, 
            is_deleted tinyint(1) NOT NULL DEFAULT '0',
            created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  		PRIMARY KEY (id),
 		KEY fk_user_marked (user_marked_id),
  		KEY fk_list (list_id),
  		CONSTRAINT fk_list FOREIGN KEY (list_id) 
  		    REFERENCES shoppinglist.sl_item_list (id) ON DELETE CASCADE ON UPDATE NO ACTION,
  		CONSTRAINT fk_user_marked FOREIGN KEY (user_marked_id) 
  		    REFERENCES shoppinglist.sl_users (id) ON DELETE CASCADE ON UPDATE NO ACTION
		) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
	`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS sl_shared_lists (
			id varchar(36) NOT NULL,
  			list_id varchar(36) NOT NULL,
  			to_user_id varchar(36) NOT NULL,
  			status tinyint(1) NOT NULL DEFAULT '0',
  			created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  			updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           	is_deleted tinyint(1) NOT NULL DEFAULT '0',
  		PRIMARY KEY (id),
  		KEY fk_user_shared_to (to_user_id),
  		KEY fk_list_id (list_id),
  		CONSTRAINT fk_list_id FOREIGN KEY (list_id) 
  		    REFERENCES shoppinglist.sl_item_list (id) ON DELETE CASCADE ON UPDATE NO ACTION,
  		CONSTRAINT fk_user_shared_to FOREIGN KEY (to_user_id) 
  		    REFERENCES shoppinglist.sl_users (id) ON DELETE CASCADE ON UPDATE NO ACTION
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DROP TRIGGER IF EXISTS set_delete_items_and_shared;`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		CREATE TRIGGER set_delete_items_and_shared 
		    AFTER UPDATE ON shoppinglist.sl_item_list 
		    FOR EACH ROW 
		    IF NEW.is_deleted=true THEN 
		        UPDATE shoppinglist.sl_item SET is_deleted=true WHERE list_id=NEW.id; 
		        UPDATE shoppinglist.sl_shared_lists SET is_deleted=true WHERE list_id=NEW.id; 
		    END IF;
	`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func (s *DataStore) CreateTransaction() (*sql.Tx, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, errors.New("Error open transaction; " + err.Error())
	}

	return tx, nil
}

func (s *DataStore) GetUsersReadRepository() readModels.UsersReadRepository {
	return readModels.NewUsersReadRepository(s.db)
}

func (s *DataStore) GetListsReadRepository() readModels.ListsReadRepository {
	return readModels.ListsReadRepository{DB: s.db}
}

func (s *DataStore) GetItemsReadRepository() readModels.ItemsReadRepository {
	return readModels.NewItemsReadRepository(s.db)
}

func (s *DataStore) GetSharesReadRepository() readModels.SharesReadRepository {
	return readModels.NewSharesReadRepository(s.db)
}

func (s *DataStore) UserProductsReadRepository() readModels.UserProductsReadRepository {
	return readModels.NewUserProductsReadRepository(s.db)
}

func (s *DataStore) GetUsersRepository(tx *sql.Tx) repositories.UsersRepository {
	if tx != nil {
		return repositories.NewUsersRepository(tx)
	}

	return repositories.NewUsersRepository(s.db)
}

func (s *DataStore) GetListsRepository(tx *sql.Tx) repositories.ListsRepository {
	return repositories.ListsRepository{DB: tx}
}

func (s *DataStore) GetItemsRepository(tx *sql.Tx) repositories.ItemsRepository {
	if tx != nil {
		return repositories.NewItemsRepository(tx)
	}

	return repositories.NewItemsRepository(s.db)
}

func (s *DataStore) GetSharesRepository(tx *sql.Tx) repositories.SharesRepository {
	if tx != nil {
		return repositories.NewSharesRepository(tx)
	}

	return repositories.NewSharesRepository(s.db)
}

func (s *DataStore) UserProductsRepository(tx *sql.Tx) repositories.UserProductsRepository {
	if tx != nil {
		return repositories.NewUserProductsRepository(tx)
	}

	return repositories.NewUserProductsRepository(s.db)
}
