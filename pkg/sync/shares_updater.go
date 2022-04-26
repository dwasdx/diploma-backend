package sync

import (
	"errors"
	"fmt"
	"shopingList/pkg"
	"shopingList/pkg/events"
	"shopingList/pkg/models"
	"shopingList/pkg/readModels"
	"shopingList/pkg/repositories"
	"time"
)

type SharesUpdater struct {
	user                 models.User
	sharesReadRepository readModels.SharesReadRepository
	sharesRepository     repositories.SharesRepository
	listsReadRepository  readModels.ListsReadRepository
	usersReadRepository  readModels.UsersReadRepository
	eventCollection      *EventCollection
}

func NewSharesUpdater(
	user models.User,
	sharesRepository repositories.SharesRepository,
	sharesReadRepository readModels.SharesReadRepository,
	listsReadRepository readModels.ListsReadRepository,
	usersReadRepository readModels.UsersReadRepository,
	eventCollection *EventCollection) SharesUpdater {

	return SharesUpdater{
		user:                 user,
		sharesRepository:     sharesRepository,
		sharesReadRepository: sharesReadRepository,
		listsReadRepository:  listsReadRepository,
		usersReadRepository:  usersReadRepository,
		eventCollection:      eventCollection}
}

func (s *SharesUpdater) Run(shares []models.ListShare) error {
	if len(shares) == 0 {
		return nil
	}

	for _, share := range shares {
		// Для шарингов обязательно нужны объекты списка
		list, err := s.listsReadRepository.GetListForIdAndOwner(share.ListID, share.OwnerID)
		if err != nil {
			if _, ok := err.(repositories.ErrNotFound); ok {
				return errors.New(fmt.Sprintf("A list from share object doesn`t found by id (%v) and owner (%v): ", share.ListID, share.OwnerID))
			}

			return errors.New("Can't get list in handleShares() by id: " + share.ListID)
		}

		if list.OwnerID != share.OwnerID {
			return errors.New("list.OwnerID != share.OwnerID for share with id: " + share.ID)
		}

		if share.OwnerID == s.user.ID {
			_, err := s.usersReadRepository.GetUser(share.ToUserID)
			if err != nil {
				if _, ok := err.(repositories.ErrNotFound); ok {
					return errors.New("User from share doesn`t found by id: " + share.ToUserID)
				}

				return errors.New("Error get user in handleShares() by id: " + share.ToUserID)
			}

			// Обработка своих шарингов
			err = s.syncOwnShare(&share, &list)
			if err != nil {
				return err
			}
		} else {
			// Обработка чужих шарингов (на пользователя)
			err = s.syncShareForUser(&share, &list)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Сохранение (создание и обновление своего шаринга)
func (s *SharesUpdater) syncOwnShare(share *models.ListShare, list *models.List) error {
	if share.ListID != list.ID {
		return errors.New("list id != share.ListID")
	}

	if share.OwnerID != s.user.ID {
		return errors.New("share.OwnerID != currentUser.ID")
	}

	if list.OwnerID != s.user.ID {
		return errors.New("list.OwnerID != currentUser.ID")
	}

	share.ReceivedAt = time.Now().UTC().Unix()

	existShare, err := s.sharesReadRepository.GetShare(share.ID, list.OwnerID)

	if err != nil {
		if err != pkg.ErrNotFoundInStorage {
			return errors.New("can't get share by id: " + share.ID + "; " + err.Error())
		}

		// Создаем новый
		if share.Status != models.ShareStatusNew {
			return errors.New("wrong status for new share with id: " + share.ID)
		}

		errCreate := s.sharesRepository.CreateShare(share)
		if errCreate != nil {
			return errors.New("error create share object: " + share.ID + "; " + errCreate.Error())
		}

		event := events.NewShareListEvent(events.ShareListEventInvite, *list, s.user, share.ToUserID)
		s.eventCollection.AddShareEvent(event)

		return nil
	}

	if share.IsEqual(*existShare) {
		return nil
	}

	// Для акцептованного шаринга уведомить получателя об удалении
	if share.Status == models.ShareStatusAccepted && !existShare.IsDeleted && share.IsDeleted {
		event := events.NewShareListEvent(events.ShareListEventDelete, *list, s.user, share.ToUserID)
		s.eventCollection.AddShareEvent(event)
	}

	// Владелец шары не может изменить ее статус. Может только удалить
	if existShare.Status != share.Status {
		share.Status = existShare.Status
	}

	err = s.sharesRepository.UpdateShare(share)
	if err != nil {
		return errors.New("can't update share")
	}

	return nil
}

// Обновить объекты шаринга, предназначенные для текущего пользователя
// Разрешено только изменения статуса и updated_at
func (s *SharesUpdater) syncShareForUser(share *models.ListShare, list *models.List) error {
	if share.ListID != list.ID {
		return errors.New("list id != share.ListID")
	}

	if share.ToUserID != s.user.ID {
		return errors.New("share.ToUserID != currentUser.id")
	}

	existShare, err := s.sharesReadRepository.GetShare(share.ID, share.OwnerID)
	if err != nil {
		return errors.New("can't get share object for me from db: " + share.ID + "; " + err.Error())
	}

	if existShare.Status != models.ShareStatusNew && share.Status == models.ShareStatusNew {
		return errors.New(fmt.Sprintf("change share status to %d forbidden; ", models.ShareStatusNew))
	}

	var oldStatus int

	oldStatus = existShare.Status
	existShare.UpdatedAt = share.UpdatedAt
	existShare.Status = share.Status
	existShare.ReceivedAt = time.Now().UTC().Unix()

	err = s.sharesRepository.UpdateShare(existShare)
	if err != nil {
		return errors.New("can't update share for me; " + err.Error())
	}

	if oldStatus != share.Status {
		if share.Status == models.ShareStatusAccepted {
			event := events.NewShareListEvent(events.ShareListEventAccept, *list, s.user, share.OwnerID)
			s.eventCollection.AddShareEvent(event)
		}

		if share.Status == models.ShareStatusRefused {
			event := events.NewShareListEvent(events.ShareListEventRefuse, *list, s.user, share.OwnerID)
			s.eventCollection.AddShareEvent(event)
		}
	}

	return nil
}
