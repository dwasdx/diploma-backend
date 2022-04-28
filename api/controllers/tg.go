package controllers

import (
	"net/http"
	"shopingList/api"
)

type TGController struct {
}

func (s *TGController) Routes() []api.Route {
	return []api.Route{
		{
			Name:   "GetRefbook",
			Method: "GET",
			Path:   "/tg",
			Func:   s.readMessage,
		},
	}
}

func (s *TGController) readMessage(w http.ResponseWriter, r *http.Request) {

}
