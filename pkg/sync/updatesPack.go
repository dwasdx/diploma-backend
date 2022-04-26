package sync

import "shopingList/pkg/models"

// Структура выдачи данных по синхронизации
type UpdatesPack struct {
	Users        []models.UserInterface `json:"users"`
	Lists        []models.List          `json:"lists"`
	Items        []models.ListItem      `json:"items"`
	Shares       []models.ListShare     `json:"shares"`
	UserProducts []models.UserProduct   `json:"user_products"`
}

func (s *UpdatesPack) GetUserIdsInObjects() []string {
	var ids []string
	var mapUserIds = make(map[string]string)

	// Собрать ID юзеров шарингов
	for _, share := range s.Shares {
		mapUserIds[share.ToUserID] = share.ToUserID
	}

	// Собрать юзеров списков
	for _, list := range s.Lists {
		mapUserIds[list.OwnerID] = list.OwnerID
	}

	for _, id := range mapUserIds {
		ids = append(ids, id)
	}

	return ids
}

func (s *UpdatesPack) IsExistList(listId string) bool {
	for _, list := range s.Lists {
		if list.ID == listId {
			return true
		}
	}

	return false
}
