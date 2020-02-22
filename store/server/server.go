package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/chadgrant/dynamodb-go-sample/store/handlers"
	"github.com/chadgrant/dynamodb-go-sample/store/repository/dynamo"
	"github.com/chadgrant/dynamodb-go-sample/store/repository/mock"
	"github.com/chadgrant/go-http-infra/infra"
	"github.com/chadgrant/go-http-infra/infra/health"
	"github.com/chadgrant/go-http-infra/infra/schema"
	"github.com/gorilla/mux"
)

type server struct {
	healthChecks health.HealthChecker
	registry     schema.Registry
	validator    schema.Validator
	dynamo       dynamodbiface.DynamoDBAPI
	handlers     *handler
	router       *mux.Router
	config       *Configuration
	errors       *log.Logger
	info         *log.Logger
}

type handler struct {
	product  *handlers.Product
	category *handlers.Category
}

func New(cfg *Configuration) (*server, error) {

	srv := &server{
		config:       cfg,
		errors:       log.New(os.Stderr, fmt.Sprintf("[%s] ERROR: ", cfg.Service.Name), log.Ldate|log.Ltime|log.Lshortfile),
		info:         log.New(os.Stdout, fmt.Sprintf("[%s] INFO: ", cfg.Service.Name), log.Ldate|log.Ltime|log.Lshortfile),
		dynamo:       dynamo.New(cfg.AWS.Region, cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, cfg.Dynamo.Endpoint),
		router:       mux.NewRouter().StrictSlash(false),
		registry:     schema.NewRegistry(),
		healthChecks: RegisterHealthChecks(cfg),
	}

	srv.router.Use(infra.Recovery)

	repo := dynamo.NewProductRepository(cfg.Dynamo.Tables.Products, srv.dynamo)
	crepo := dynamo.NewCategoryRepository(cfg.Dynamo.Tables.Categories, srv.dynamo)
	if cfg.UseMocks {
		crepo = mock.NewCategoryRepository("Hats", "Shirts", "Pants", "Shoes", "Ties", "Belts", "Socks", "Accessory")
		repo = mock.NewProductRepository(crepo, 100)
	}

	srv.handlers = &handler{
		product:  handlers.NewProduct(srv.errors, repo),
		category: handlers.NewCategory(srv.errors, crepo),
	}
	if err := srv.registerRoutes(); err != nil {
		return nil, err
	}

	return srv, nil
}

func (s *server) Serve(address string) error {
	s.info.Printf("Started serving at %s\n", address)
	s.errors.Fatal(http.ListenAndServe(address, s.router))
	return nil
}
