package handlers

import (
	"log"
	"net/http"

	"github.com/chadgrant/dynamodb-go-sample/store/repository"
)

type Category struct {
	errors *log.Logger
	repo   repository.CategoryRepository
}

func NewCategory(errors *log.Logger, repo repository.CategoryRepository) *Category {
	return &Category{errors, repo}
}

func (h *Category) GetAll(w http.ResponseWriter, r *http.Request) {
	cats, err := h.repo.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.errors.Printf("getting categories: %v\n", err)
		return
	}

	returnJSON(w, r, cats)
}
