package handlers

import (
	"log"
	"net/http"

	"github.com/chadgrant/dynamodb-go-sample/store/service"
)

type Category struct {
	errors *log.Logger
	svc    service.Service
}

func NewCategory(errors *log.Logger, svc service.Service) *Category {
	return &Category{errors, svc}
}

func (h *Category) GetAll(w http.ResponseWriter, r *http.Request) {
	cats, err := h.svc.Categories()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.errors.Printf("getting categories: %v\n", err)
		return
	}

	returnJSON(w, r, cats)
}
