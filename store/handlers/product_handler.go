package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func (h *ProductHandler) Categories(w http.ResponseWriter, r *http.Request) {
	//faked
	var categories = []string{"Hats", "Shirts", "Pants", "Shoes", "Ties", "Belts", "Socks", "Accessory"}

	json.NewEncoder(w).Encode(categories)
}

func (h *ProductHandler) GetPaged(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cat := vars["category"]

	last := param(r, "last", "")
	lastprice, _ := strconv.ParseFloat(param(r, "lastprice", "0"), 2)

	products, _, err := h.repo.GetPaged(cat, 25, last, lastprice)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(products)
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
	vars := mux.Vars(r)
	cat := vars["category"]

	var p store.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.repo.Upsert(cat, &p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("product/%s", p.ID))
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	cat := vars["category"]

	var p store.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p.ID = id

	if err := h.repo.Upsert(cat, &p); err != nil {
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

func param(r *http.Request, key, defaultValue string) string {
	v, ok := r.URL.Query()[key]
	if ok && len(v) > 0 && len(v[0]) > 0 {
		return v[0]
	}
	return defaultValue
}
