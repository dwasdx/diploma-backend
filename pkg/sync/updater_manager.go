package sync

import (
	"github.com/pkg/errors"
	"log"
	"shopingList/pkg/events"
	"shopingList/pkg/models"
	"shopingList/pkg/repositories"
	"shopingList/store"
)

type UpdaterManager struct {
	dataService     store.DataService
	user            models.User
	ChanGoodsChange chan events.GoodsChangeEvent
	ChanShareChange chan events.ShareListEvent
	eventCollection EventCollection
}

func NewUpdater(dataService store.DataService, user models.User) *UpdaterManager {
	return &UpdaterManager{
		dataService: dataService,
		user:        user}
}

func (s *UpdaterManager) RunUpdate(users []models.User, lists []models.List, shares []models.ListShare, items []models.ListItem, userProducts []models.UserProduct) error {
	tx, err := s.dataService.CreateTransaction()
	if err != nil {
		return errors.New("Error open transaction; " + err.Error())
	}

	defer tx.Rollback()

	usersReadRepository := s.dataService.GetUsersReadRepository()
	usersRepository := s.dataService.GetUsersRepository(tx)
	itemsRepository := s.dataService.GetItemsRepository(tx)
	listsReadRepository := s.dataService.GetListsReadRepository()
	listsRepository := s.dataService.GetListsRepository(tx)
	sharesReadRepository := s.dataService.GetSharesReadRepository()
	sharesRepository := s.dataService.GetSharesRepository(tx)
	userProductsRepository := s.dataService.UserProductsRepository(tx)
	userProductsReadRepository := s.dataService.UserProductsReadRepository()
	listsCollection := NewListsCollection(&listsReadRepository)

	// Обновить пользователей
	err = s.handleUsers(users, &usersRepository)
	if err != nil {
		return errors.New("Error update users; " + err.Error())
	}

	// Обновить списки
	listsUpdater := NewUpdaterList(s.user.ID, &listsRepository, listsCollection)
	err = listsUpdater.Run(lists)
	if err != nil {
		return errors.New("Error update lists; " + err.Error())
	}

	// Обновить шаринги
	sharesUpdater := NewSharesUpdater(s.user, sharesRepository, sharesReadRepository, listsReadRepository, usersReadRepository, &s.eventCollection)
	err = sharesUpdater.Run(shares)
	if err != nil {
		return errors.New("Error update shares; " + err.Error())
	}

	// Обновить товары
	itemsUpdater := NewItemsUpdater(s.user, listsCollection, itemsRepository, listsReadRepository,
		sharesReadRepository, usersReadRepository, &s.eventCollection)

	err = itemsUpdater.Run(items, lists)
	if err != nil {
		return errors.New("Error update items; " + err.Error())
	}

	// Обновить товары пользователя
	userProductsUpdater := UserProductsUpdater(s.user, userProductsRepository, userProductsReadRepository, &s.eventCollection)
	err = userProductsUpdater.Run(userProducts)
	if err != nil {
		return errors.New("Error update userProducts; " + err.Error())
	}

	if err = tx.Commit(); err != nil {
		log.Fatal(errors.Wrap(err, "Error commit"))
	}

	s.sendEvents()

	return nil
}

func (s *UpdaterManager) sendEvents() {
	for _, event := range s.eventCollection.GetShareEvents() {
		s.ChanShareChange <- event
	}

	for _, event := range s.eventCollection.GetGoodEvents() {
		s.ChanGoodsChange <- event
	}
}

func (s *UpdaterManager) handleUsers(users []models.User, usersRepository *repositories.UsersRepository) error {
	if len(users) == 0 {
		return nil
	}

	for _, u := range users {
		// Обновляем только пользователя владельца
		// Объекты других пользователей могут приходить, но мы их игнорируем
		if u.ID != s.user.ID {
			continue
		}

		if u.Phone != s.user.Phone {
			return errors.New("you may not to change phone here")
		}

		err := usersRepository.UpdateUser(&u)
		if err != nil {
			return err
		}
	}

	return nil
}
