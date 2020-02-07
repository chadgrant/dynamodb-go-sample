package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chadgrant/dynamodb-go-sample/store"
	"github.com/chadgrant/dynamodb-go-sample/store/repository/mock"
	"github.com/gorilla/mux"
)

func TestProductHandler(t *testing.T) {
	c := mock.NewCategoryRepository("Hats", "Shirts", "Pants", "Shoes", "Ties", "Belts", "Socks", "Accessory")
	repo := mock.NewProductRepository(c, 100)

	h := NewProductHandler(repo)

	t.Run("Add", func(t *testing.T) {
		testAdd(h, t)
	})

	t.Run("GetPaged", func(t *testing.T) {
		testGetPaged(h, t)
	})

	t.Run("Get", func(t *testing.T) {
		prds, err := repo.GetPaged("hats", 25, "", float64(0))
		if err != nil {
			t.Fatal(err)
		}
		if len(prds) == 0 {
			t.Fatalf("no products returned")
		}
		for _, p := range prds {
			testGet(p.ID, h, t)
		}
	})

	t.Run("UpdateProduct", func(t *testing.T) {
		prds, err := repo.GetPaged("hats", 25, "", float64(0))
		if err != nil {
			t.Fatal(err)
		}
		if len(prds) == 0 {
			t.Fatalf("no products returned")
		}
		testUpsert(prds[0], h, t)
	})

	t.Run("DeleteProduct", func(t *testing.T) {
		prds, err := repo.GetPaged("hats", 25, "", float64(0))
		if err != nil {
			t.Fatal(err)
		}
		if len(prds) == 0 {
			t.Fatalf("no products returned")
		}
		testDelete(prds[0].ID, h, t)
	})
}

func testAdd(handler *ProductHandler, t *testing.T) {
	b := []byte("{ \"name\":\"created from web test\", \"category\": \"hats\",  \"description\": \"nice product from web test\", \"price\": 5.77 }")
	r, _ := http.NewRequest(http.MethodPost, "products/hats", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	handler.Add(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected created respose got %d", w.Code)
	}
}

func testGetPaged(handler *ProductHandler, t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "products/hats", nil)
	w := httptest.NewRecorder()
	m := map[string]string{
		"category": "hats",
	}
	r = mux.SetURLVars(r, m)

	handler.GetPaged(w, r)

	js, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	var p pagedProducts
	if err := json.Unmarshal(js, &p); err != nil {
		t.Fatal(err)
	}

	if len(p.Results) == 0 {
		t.Errorf("got no products back")
	}
}

func testGet(id string, handler *ProductHandler, t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "product/"+id, nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"id": id,
	}

	r = mux.SetURLVars(r, vars)

	handler.Get(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 status got %d", w.Code)
	}

	js, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
		return
	}

	var p store.Product
	if err := json.Unmarshal(js, &p); err != nil {
		t.Fatal(err)
	}

	if p.ID != id {
		t.Errorf("unexpected product returned. expected %s got %s", id, p.ID)
	}
}

func testUpsert(product *store.Product, handler *ProductHandler, t *testing.T) {
	copy := product
	copy.Name = product.Name + " Updated"

	b, err := json.Marshal(copy)
	if err != nil {
		t.Fatal(err)
	}

	r, _ := http.NewRequest(http.MethodPut, "product/"+product.ID, bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	vars := map[string]string{
		"id": product.ID,
	}

	r = mux.SetURLVars(r, vars)

	handler.Upsert(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("unexpected status, expected %d got %d", http.StatusAccepted, w.Code)
	}

	r, _ = http.NewRequest(http.MethodGet, "product/"+copy.ID, nil)
	w = httptest.NewRecorder()

	r = mux.SetURLVars(r, vars)

	handler.Get(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 status got %d", w.Code)
	}

	js, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	var p store.Product
	if err := json.Unmarshal(js, &p); err != nil {
		t.Fatal(err)
	}

	if p.Name != copy.Name {
		t.Errorf("product name was not updated: %s", p.Name)
	}
}

func testDelete(productID string, handler *ProductHandler, t *testing.T) {
	r, _ := http.NewRequest(http.MethodDelete, "product/"+productID, nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"id": productID,
	}

	r = mux.SetURLVars(r, vars)

	handler.Delete(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("unexpected status, expected %d got %d", http.StatusAccepted, w.Code)
		return
	}

	r, _ = http.NewRequest(http.MethodGet, "product/"+productID, nil)
	w = httptest.NewRecorder()
	r = mux.SetURLVars(r, vars)

	handler.Get(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected not found but got %d", w.Code)
	}
}
