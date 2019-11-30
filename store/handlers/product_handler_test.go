package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chadgrant/dynamodb-go-sample/store"

	"github.com/chadgrant/dynamodb-go-sample/store/repository"
	"github.com/chadgrant/dynamodb-go-sample/store/repository/mock"
	"github.com/gorilla/mux"
)

func TestProductHandler(t *testing.T) {
	repo, err := setup()
	if err != nil {
		t.Fatal(err)
		return
	}

	h := NewProductHandler(repo)

	t.Run("AddProduct", func(t *testing.T) {
		testAddProduct(h, t)
	})

	t.Run("Get", func(t *testing.T) {
		testGetProducts(h, t)
	})

	t.Run("GetProduct", func(t *testing.T) {
		prds, _, err := repo.GetPaged("hats", 25, "", float64(0))
		if err != nil {
			t.Error(err)
			return
		}
		for _, p := range prds {
			testGetProduct(p.ID, h, t)
		}
	})

	t.Run("UpdateProduct", func(t *testing.T) {
		prds, _, err := repo.GetPaged("hats", 25, "", float64(0))
		if err != nil {
			t.Error(err)
			return
		}
		testUpdateProduct(prds[0], h, t)
	})

	t.Run("DeleteProduct", func(t *testing.T) {
		prds, _, err := repo.GetPaged("hats", 25, "", float64(0))
		if err != nil {
			t.Error(err)
			return
		}
		testDeleteProduct(prds[0].ID, h, t)
	})
}

func setup() (repository.ProductRepository, error) {
	repo := mock.NewProductRepository()

	populator := repository.NewPopulator(repo)

	// if err := populator.CreateProducts(100); err != nil {
	// 	return nil,err
	// }

	// if err := populator.ExportProducts("../../data/products.json"); err != nil {
	// 	return nil,err
	// }

	if err := populator.LoadProducts("../../data/products.json"); err != nil {
		return nil, err
	}

	return repo, nil
}

func testAddProduct(handler *ProductHandler, t *testing.T) {
	b := []byte("{ \"name\":\"created from web test\", \"description\": \"nice product from web test\", \"price\": 5.77 }")
	r, _ := http.NewRequest(http.MethodPost, "product/hats", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	handler.Add(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected no content response got %d", w.Code)
		return
	}
}

func testGetProducts(handler *ProductHandler, t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "product/hats", nil)
	w := httptest.NewRecorder()

	handler.GetPaged(w, r)

	js, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
		return
	}

	var p []store.Product
	if err := json.Unmarshal(js, &p); err != nil {
		t.Error(err)
		return
	}

	if len(p) == 0 {
		t.Errorf("got no products back")
	}
}

func testGetProduct(id string, handler *ProductHandler, t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "product/hats/"+id, nil)
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
		t.Error(err)
		return
	}

	if p.ID != id {
		t.Errorf("unexpected product returned. expected %s got %s", id, p.ID)
	}
}

func testUpdateProduct(product *store.Product, handler *ProductHandler, t *testing.T) {
	u := product
	u.Name = product.Name + " Updated"

	b, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
		return
	}

	r, _ := http.NewRequest(http.MethodPut, "product/hats/"+product.ID, bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	vars := map[string]string{
		"id": product.ID,
	}

	r = mux.SetURLVars(r, vars)

	handler.Upsert(w, r)

	if w.Code != http.StatusAccepted {
		t.Errorf("unexpected status, expected %d got %d", http.StatusAccepted, w.Code)
		return
	}

	r, _ = http.NewRequest(http.MethodGet, "product/hats/"+u.ID, nil)
	w = httptest.NewRecorder()

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
		t.Error(err)
		return
	}

	if p.Name != u.Name {
		t.Errorf("product name was not updated: %s", p.Name)
	}
}

func testDeleteProduct(productID string, handler *ProductHandler, t *testing.T) {
	r, _ := http.NewRequest(http.MethodDelete, "product/hats/"+productID, nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"id": productID,
	}

	r = mux.SetURLVars(r, vars)

	handler.Delete(w, r)

	if w.Code != http.StatusAccepted {
		t.Errorf("unexpected status, expected %d got %d", http.StatusAccepted, w.Code)
		return
	}

	r, _ = http.NewRequest(http.MethodGet, "product/hats/"+productID, nil)
	w = httptest.NewRecorder()
	r = mux.SetURLVars(r, vars)

	handler.Get(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected not found but got %d", w.Code)
	}
}
