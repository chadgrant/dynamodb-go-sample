package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/chadgrant/dynamodb-go-sample/store/repository"
	"github.com/chadgrant/dynamodb-go-sample/store/repository/mock"

	"github.com/chadgrant/dynamodb-go-sample/store/handlers"
	"github.com/chadgrant/go/http/infra"
	"github.com/gorilla/mux"
)

func main() {
	host := *flag.String("host", "0.0.0.0", "default binding 0.0.0.0")
	port := *flag.Int("port", 8080, "default port 8080")
	flag.Parse()

	infra.Handle()

	http.Handle("/", http.FileServer(http.Dir("docs")))

	repo := mock.NewProductRepository()
	pop := repository.NewPopulator(repo)

	if err := pop.LoadProducts("data/products.json"); err != nil {
		log.Fatalf("loading products %v", err)
		return
	}

	r := mux.NewRouter()

	ph := handlers.NewProductHandler(repo)
	r.HandleFunc("/product", ph.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/product", ph.Add).Methods(http.MethodPost)
	r.HandleFunc("/product/{id}", ph.Upsert).Methods(http.MethodPut)
	r.HandleFunc("/product/{id}", ph.Get).Methods(http.MethodGet)
	r.HandleFunc("/product/{id}", ph.Delete).Methods(http.MethodDelete)

	log.Printf("Started, serving at %s:%d\n", host, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), r))
}
