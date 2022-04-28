package sync

import (
	models "shopingList/pkg/models"
)

// ShoppingListUpdatesRequest - query model for shopping list update request
type ShoppingListUpdatesRequest struct {
	Date int64 `json:"date"`
}

// ShoppingListUpdates - complex object of all shopping list related entities
type ShoppingListUpdates struct {
	Users        []models.User        `json:"users"`
	Lists        []models.List        `json:"lists"`
	Items        []models.ListItem    `json:"items"`
	Shares       []models.ListShare   `json:"shares"`
	UserProducts []models.UserProduct `json:"user_products"`
}

func (s *ShoppingListUpdates) Validate() (bool, []string) {
	var errs []string

	for _, user := range s.Users {
		_, err := user.Validate()
		if err != nil {
			errs = append(errs, "error in user "+user.ID+"; "+err.Error())
		}
	}

	for _, item := range s.Items {
		_, err := item.Validate()
		if err != nil {
			errs = append(errs, "error in item "+item.ID+"; "+err.Error())
		}
	}

	for _, list := range s.Lists {
		_, err := list.Validate()
		if err != nil {
			errs = append(errs, "error in list "+list.ID+"; "+err.Error())
		}
	}

	for _, share := range s.Shares {
		_, err := share.Validate()
		if err != nil {
			errs = append(errs, "error in share "+share.ID+"; "+err.Error())
		}
	}

	for _, userProduct := range s.UserProducts {
		_, err := userProduct.Validate()
		if err != nil {
			errs = append(errs, "error in userProduct "+userProduct.ID+"; "+err.Error())
		}
	}

	if len(errs) > 0 {
		return false, errs
	}

	return true, nil
}
