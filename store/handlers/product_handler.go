package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chadgrant/dynamodb-go-sample/store"
	"github.com/chadgrant/dynamodb-go-sample/store/repository"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo}
}

func (h *ProductHandler) GetPaged(w http.ResponseWriter, r *http.Request) {
	// cat := r.URL.Query()["category"][0]
	// last := r.URL.Query()["last"][0]
	// lastprice, _ := strconv.ParseFloat(r.URL.Query()["lastprice"][0], 2)

	// products, _, err := h.repo.GetPaged(cat, 25, last, lastprice)
	// if err != nil {
	// 	return
	// }
	// json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	p, err := h.repo.Get(id)
	if err != nil {
		return
	}

	if p == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) Add(w http.ResponseWriter, r *http.Request) {

	var p store.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.repo.Upsert(&p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("product/%s", p.ID))
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var p store.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p.ID = id

	if err := h.repo.Upsert(&p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
