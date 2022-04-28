package controllers

import (
	"net/http"
	"shopingList/api"
	"shopingList/pkg/repositories"
)

type RefbookController struct {
	categoriesRepository repositories.RefbookCategoriesRepository
	productsRepository   repositories.RefbookProductsRepository
}

func NewRefbookController(categoriesRepository repositories.RefbookCategoriesRepository,
	productsRepository repositories.RefbookProductsRepository) *RefbookController {
	return &RefbookController{categoriesRepository: categoriesRepository, productsRepository: productsRepository}
}

func (s *RefbookController) Routes() []api.Route {
	return []api.Route{
		{
			Name:   "GetRefbook",
			Method: "GET",
			Path:   "/refbook",
			Func:   s.getRefbook,
		},
	}
}

func (s *RefbookController) getRefbook(w http.ResponseWriter, r *http.Request) {
	categories, err := s.categoriesRepository.GetAll()
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can`t get categories", api.ErrInternal)
		return
	}

	products, err := s.productsRepository.GetAll()
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get products", api.ErrInternal)
		return
	}

	data := make(map[string]interface{})
	data["categories"] = categories
	data["products"] = products

	api.SendDataJSON(w, r, http.StatusOK, data)
}
