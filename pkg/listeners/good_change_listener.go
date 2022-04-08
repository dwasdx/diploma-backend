package listeners

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"shopingList/pkg/events"
	"shopingList/pkg/models"
	"shopingList/pkg/repositories"
	"shopingList/pkg/services"
)

type GoodChangeListener struct {
	Repository  repositories.NotificationsRepository
	PushChannel chan services.PushNotificationMessage
}

func (s *GoodChangeListener) Run(channel chan events.GoodsChangeEvent) {
	for event := range channel {
		log.Info("Receive message for GoodChangeListener")
		err := s.Handle(&event)
		if err != nil {
			log.Error("Error handle event", err)
		}
	}
}

func (s *GoodChangeListener) Handle(event interface{}) error {
	model, ok := event.(*events.GoodsChangeEvent)

	if !ok {
		log.Fatalf("event must be a type not GoodsChangeEvent, event is type: %#v/n", event)
	}

	var message string

	typeEvent := int(model.TypeNotification())
	item := model.Item()
	user := model.User()
	listName := model.ListName()

	switch typeEvent {
	case models.NotificationTypeGoodsCreate:
		message = fmt.Sprintf("Пользователь %d добавил товар \"%s\" в список \"%s\"", user.Phone, item.Name, listName)
	case models.NotificationTypeGoodsChange:
		message = fmt.Sprintf("Пользователь %d изменил товар \"%s\" из списка \"%s\"", user.Phone, item.Name, listName)
	case models.NotificationTypeGoodsCheck:
		message = fmt.Sprintf("Пользователь %d отметил товар \"%s\" из списка \"%s\"", user.Phone, item.Name, listName)
	case models.NotificationTypeGoodsUncheck:
		message = fmt.Sprintf("Пользователь %d снял отметку с товара \"%s\" из списка \"%s\"", user.Phone, item.Name, listName)
	case models.NotificationTypeGoodsDelete:
		message = fmt.Sprintf("Пользователь %d удалил товар \"%s\" из списка \"%s\"", user.Phone, item.Name, listName)
	}

	form := models.NotificationCreateForm{
		TypeNotification: model.TypeNotification(),
		Message:          message,
		UserId:           user.ID,
		UserPhone:        user.Phone,
		ListId:           model.Item().ListID,
		ItemId:           models.NullString{String: model.Item().ID, Valid: true},
	}

	// Создать уведомления
	for _, targetUserId := range model.TargetUserIds() {
		form.TargetUserId = targetUserId
		err := s.Repository.Create(&form)

		if err != nil {
			log.Printf("Error create notification" + err.Error())
			return err
		}

	}

	if s.PushChannel == nil {
		log.Println("[ERROR] PushChannel is nil")
	} else {
		pushMessage := services.PushNotificationMessage{Notification: form, TargetUserIds: model.TargetUserIds()}
		s.PushChannel <- pushMessage
	}

	return nil
}
