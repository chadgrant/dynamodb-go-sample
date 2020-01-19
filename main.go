package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chadgrant/dynamodb-go-sample/store/handlers"
	"github.com/chadgrant/dynamodb-go-sample/store/repository"
	"github.com/chadgrant/dynamodb-go-sample/store/repository/dynamo"
	"github.com/chadgrant/go-http-infra/gorilla"
	"github.com/chadgrant/go-http-infra/infra"
	"github.com/chadgrant/go-http-infra/infra/health"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	host := *flag.String("host", infra.GetEnvVar("SVC_HOST", "0.0.0.0"), "default binding 0.0.0.0")
	port := *flag.Int("port", infra.GetEnvVarInt("SVC_PORT", 8080), "default port 8080")
	mock := *flag.Bool("mock", infra.GetEnvVarBool("SVC_MOCK_DATA", false), "use mock database")
	region := *flag.String("region", infra.GetEnvVar("AWS_REGION", "us-east-1"), "aws region")
	//prefer using IAM roles, but local dynamo demands creds....
	accessKey := *flag.String("accessKey", infra.GetEnvVar("AWS_ACCESS_KEY_ID", "key"), "aws access key")
	keySecret := *flag.String("keySecret", infra.GetEnvVar("AWS_SECRET_ACCESS_KEY", "secret"), "aws access key secret")
	sessionToken := *flag.String("sessionToken", infra.GetEnvVar("AWS_SESSION_TOKEN", ""), "aws session token")
	endpoint := *flag.String("endpoint", infra.GetEnvVar("DYNAMO_ENDPOINT", "http://localhost:8000"), "dynamo endpoint url")
	table := *flag.String("table", infra.GetEnvVar("DYNAMO_TABLE", "products"), "products dynamo table name")
	ctable := *flag.String("ctable", infra.GetEnvVar("DYNAMO_TABLE", "categories"), "categories dynamo table name")
	flag.Parse()

	dyn := dynamodb.New(session.Must(session.NewSession()), &aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, keySecret, sessionToken),
		Endpoint:    aws.String(endpoint),
	})

	var repo repository.ProductRepository
	var crepo repository.CategoryRepository
	repo = dynamo.NewProductRepository(table, dyn)
	crepo = dynamo.NewCategoryRepository(ctable, dyn)
	if mock {
		crepo = repository.NewMockCategoryRepository("Hats", "Shirts", "Pants", "Shoes", "Ties", "Belts", "Socks", "Accessory")
		repo = repository.NewMockProductRepository(crepo, 100)
	}

	checker := health.NewHealthChecker()
	checker.AddReadiness("dynamo", dynamo.Health(dyn, "products", time.Second*1))

	hc := health.NewHandler(checker)
	r := mux.NewRouter()
	gorilla.Handle(r, hc)
	r.Use(infra.Recovery)

	ph := handlers.NewProductHandler(repo)
	ch := handlers.NewCategoryHandler(crepo)

	r.HandleFunc("/category", ch.GetAll).Methods(http.MethodGet)

	r.HandleFunc("/product/{category:[A-Za-z]+}", ph.GetPaged).Methods(http.MethodGet)
	r.HandleFunc("/product/", ph.Add).Methods(http.MethodPost)
	r.HandleFunc("/product/{id:[a-z0-9\\-]{36}}", ph.Upsert).Methods(http.MethodPut)
	r.HandleFunc("/product/{id:[a-z0-9\\-]{36}}", ph.Get).Methods(http.MethodGet)
	r.HandleFunc("/product/{id:[a-z0-9\\-]{36}}", ph.Delete).Methods(http.MethodDelete)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./docs/")))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		ExposedHeaders:   []string{"Location"},
		MaxAge:           86400,
	})

	log.Printf("Started, serving at %s:%d\n", host, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), c.Handler(r)))
}
