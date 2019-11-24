package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/chadgrant/go/http/infra"
)


func main() {
	host := *flag.String("host", "0.0.0.0", "default binding 0.0.0.0")
	port := *flag.Int("port", 8080, "default port 8080")
	flag.Parse()
	
	infra.Handle()

	http.Handle("/", http.FileServer(http.Dir("docs")))

	log.Printf("Started, serving at %s:%d\n", host, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))
}