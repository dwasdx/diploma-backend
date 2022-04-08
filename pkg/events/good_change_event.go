package events

import (
	"fmt"
	"shopingList/pkg/models"
)

type GoodsChangeEvent struct {
	typeNotification models.NotificationType
	item             models.ListItem
	listName         string
	user             *models.User
	targetUserIds    []string
}

func NewGoodChangeEvent(
	typeEvent models.NotificationType,
	item models.ListItem,
	user *models.User,
	listName string,
	targetUserIds []string) GoodsChangeEvent {
	checkIds := make(map[string]bool)
	uniqueIds := make([]string, 0)

	for _, userId := range targetUserIds {
		if _, ok := checkIds[userId]; !ok {
			uniqueIds = append(uniqueIds, userId)
		}

		checkIds[userId] = true
	}

	return GoodsChangeEvent{typeNotification: typeEvent, item: item, user: user, listName: listName, targetUserIds: uniqueIds}
}

func (s *GoodsChangeEvent) GetEventType() string {
	return fmt.Sprint("%s")
}

func (s *GoodsChangeEvent) TypeNotification() models.NotificationType {
	return s.typeNotification
}

func (s *GoodsChangeEvent) Item() *models.ListItem {
	return &s.item
}

func (s *GoodsChangeEvent) ListName() string {
	return s.listName
}

func (s *GoodsChangeEvent) User() *models.User {
	return s.user
}

func (s *GoodsChangeEvent) TargetUserIds() []string {
	return s.targetUserIds
}
