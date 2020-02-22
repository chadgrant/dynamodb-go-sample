package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/chadgrant/dynamodb-go-sample/store"
	"github.com/chadgrant/dynamodb-go-sample/store/repository"
	"github.com/chadgrant/go-http-infra/infra"
)

type Product struct {
	errors *log.Logger
	repo   repository.ProductRepository
}

type pagedProducts struct {
	Results []*store.Product `json:"results"`
	Next    string           `json:"next,omitempty"`
}

func NewProduct(errors *log.Logger, repo repository.ProductRepository) *Product {
	return &Product{errors, repo}
}

func (h *Product) GetPaged(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cat := vars["category"]

	last := param(r, "last", "")
	lastprice, _ := strconv.ParseFloat(param(r, "lastprice", "0"), 2)

	products, err := h.repo.GetPaged(cat, 25, last, lastprice)
	if err != nil {
		h.errors.Printf("getting products: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
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

func (h *Product) Get(w http.ResponseWriter, r *http.Request) {
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

func (h *Product) Add(w http.ResponseWriter, r *http.Request) {
	var p store.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.errors.Printf("adding product (new random): %v\n", err)
		return
	}

	p.ID = id.String()

	if err := h.repo.Upsert(&p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.errors.Printf("adding product: %v\n", err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/product/%s", p.ID))
	w.WriteHeader(http.StatusCreated)
}

func (h *Product) Upsert(w http.ResponseWriter, r *http.Request) {
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
		h.errors.Printf("updating product: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Product) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.errors.Printf("deleting product: %v\n", err)
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
