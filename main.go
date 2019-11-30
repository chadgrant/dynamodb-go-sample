package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/chadgrant/dynamodb-go-sample/store/handlers"
	"github.com/chadgrant/dynamodb-go-sample/store/repository"
	"github.com/chadgrant/go/http/infra"
	"github.com/gorilla/mux"
)

func main() {
	host := *flag.String("host", "0.0.0.0", "default binding 0.0.0.0")
	port := *flag.Int("port", 8080, "default port 8080")
	flag.Parse()

	//repo := dynamo.NewProductRepository()
	repo := repository.NewMockProductRepository()

	pop := repository.NewPopulator(repo)
	if err := pop.Load("data/products.json"); err != nil {
		log.Fatalf("loading products %v", err)
		return
	}

	r := mux.NewRouter()
	//r.Use(contentType)

	infra.HandleGorilla(r)

	//fs := http.FileServer(http.Dir("docs"))
	//r.Handle("/docs/", http.StripPrefix("/docs/", fs))

	ph := handlers.NewProductHandler(repo)

	r.HandleFunc("/category", ph.Categories).Methods(http.MethodGet)

	r.HandleFunc("/product/{category}", ph.GetPaged).Methods(http.MethodGet)
	r.HandleFunc("/product/{category}", ph.Add).Methods(http.MethodPost)
	r.HandleFunc("/product/{category}/{id}", ph.Upsert).Methods(http.MethodPut)
	r.HandleFunc("/product/{category}/{id}", ph.Get).Methods(http.MethodGet)
	r.HandleFunc("/product/{category}/{id}", ph.Delete).Methods(http.MethodDelete)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./docs/")))

	log.Printf("Started, serving at %s:%d\n", host, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), r))
}

func contentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
