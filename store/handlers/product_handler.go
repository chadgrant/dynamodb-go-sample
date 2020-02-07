package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/chadgrant/dynamodb-go-sample/store"
	"github.com/chadgrant/dynamodb-go-sample/store/repository"
	"github.com/chadgrant/go-http-infra/infra"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

type pagedProducts struct {
	Results []*store.Product `json:"results"`
	Next    string           `json:"next,omitempty"`
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo}
}

func (h *ProductHandler) GetPaged(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cat := vars["category"]

	last := param(r, "last", "")
	lastprice, _ := strconv.ParseFloat(param(r, "lastprice", "0"), 2)

	products, err := h.repo.GetPaged(cat, 25, last, lastprice)
	if err != nil {
		return
	}

	next := ""
	if len(products) > 0 {
		p := products[len(products)-1]
		next = fmt.Sprintf("/products/%s?last=%s&lastprice=%.2f", cat, p.ID, p.Price)
	}

	returnJSON(w, r, &pagedProducts{
		Results: products,
		Next:    next,
	})
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

	returnJSON(w, r, p)
}

func (h *ProductHandler) Add(w http.ResponseWriter, r *http.Request) {
	var p store.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p.ID = id.String()

	if err := h.repo.Upsert(&p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/product/%s", p.ID))
	w.WriteHeader(http.StatusCreated)
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

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func param(r *http.Request, key, defaultValue string) string {
	v, ok := r.URL.Query()[key]
	if ok && len(v) > 0 && len(v[0]) > 0 {
		return v[0]
	}
	return defaultValue
}

func returnJSON(w http.ResponseWriter, r *http.Request, o interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(o); err != nil {
		infra.Error(w, r, http.StatusInternalServerError, err)
	}
}
