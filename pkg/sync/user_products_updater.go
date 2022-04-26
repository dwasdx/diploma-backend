package sync

import (
	"errors"
	"shopingList/pkg/models"
	"shopingList/pkg/readModels"
	"shopingList/pkg/repositories"
	"time"
)

type userProductsUpdater struct {
	user                       models.User
	userProductsRepository     repositories.UserProductsRepository
	userProductsReadRepository readModels.UserProductsReadRepository
	usersReadRepository        readModels.UsersReadRepository
	notificationService        NotificationsCreateService
	eventCollection            *EventCollection
}

func UserProductsUpdater(
	user models.User,
	userProductsRepository repositories.UserProductsRepository,
	UserProductsReadRepository readModels.UserProductsReadRepository,
	eventCollection *EventCollection) *userProductsUpdater {
	return &userProductsUpdater{
		user:                       user,
		userProductsRepository:     userProductsRepository,
		userProductsReadRepository: UserProductsReadRepository,
		eventCollection:            eventCollection,
		notificationService:        NotificationsCreateService{}}
}

func (s *userProductsUpdater) getUserId() string {
	return s.user.ID
}

func (s *userProductsUpdater) Run(userProducts []models.UserProduct) error {
	if len(userProducts) == 0 {
		return nil
	}

	for _, item := range userProducts {
		item.ReceivedAt = time.Now().UTC().Unix()

		existItem, err := s.userProductsRepository.GetOneById(item.ID)
		if err != nil {
			if _, ok := err.(repositories.ErrNotFound); ok {
				err = s.userProductsRepository.Create(&item)
				if err != nil {
					return errors.New("Can`t create item; " + err.Error())
				}
			}

			continue
		}

		if existItem.IsEqual(&item) {
			continue
		}

		err = s.userProductsRepository.Update(&item)
		if err != nil {
			return errors.New("Can`t update item; " + err.Error())
		}

	}

	return nil
}
