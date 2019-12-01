package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chadgrant/dynamodb-go-sample/store"
	"github.com/chadgrant/dynamodb-go-sample/store/repository"
	"github.com/chadgrant/go/http/infra"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

type pagedProducts struct {
	Results []*store.Product `json:"results"`
	Next    string           `json:"next,omitempty"`
	Total   int64            `json:"total"`
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo}
}

func (h *ProductHandler) Categories(w http.ResponseWriter, r *http.Request) {
	//faked
	var categories = []string{"Hats", "Shirts", "Pants", "Shoes", "Ties", "Belts", "Socks", "Accessory"}

	returnJson(w, r, categories)
}

func (h *ProductHandler) GetPaged(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cat := vars["category"]

	last := param(r, "last", "")
	lastprice, _ := strconv.ParseFloat(param(r, "lastprice", "0"), 2)

	products, total, err := h.repo.GetPaged(cat, 25, last, lastprice)
	if err != nil {
		return
	}

	next := ""
	if len(products) > 0 {
		p := products[len(products)-1]
		next = fmt.Sprintf("/product/%s/?last=%s&lastPrice=%.2f", cat, p.ID, p.Price)
	}

	returnJson(w, r, &pagedProducts{
		Results: products,
		Total:   total,
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

	returnJson(w, r, p)
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

func returnJson(w http.ResponseWriter, r *http.Request, o interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(o); err != nil {
		infra.Error(w, r, http.StatusInternalServerError, err)
	}
}
