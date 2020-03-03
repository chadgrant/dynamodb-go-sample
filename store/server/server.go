package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chadgrant/dynamodb-go-sample/store/handlers"
	"github.com/chadgrant/dynamodb-go-sample/store/metrics"
	"github.com/chadgrant/dynamodb-go-sample/store/repository"
	"github.com/chadgrant/dynamodb-go-sample/store/repository/dynamo"
	"github.com/chadgrant/dynamodb-go-sample/store/repository/mock"
	"github.com/chadgrant/dynamodb-go-sample/store/service"
	"github.com/chadgrant/go-http-infra/infra"
	"github.com/chadgrant/go-http-infra/infra/health"
	"github.com/chadgrant/go-http-infra/infra/schema"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/push"
)

type Server struct {
	healthChecks health.HealthChecker
	registry     schema.Registry
	validator    schema.Validator
	dynamo       dynamodbiface.DynamoDBAPI
	handlers     *handler
	router       *mux.Router
	config       *Configuration
	errors       *log.Logger
	info         *log.Logger
	server       *http.Server
}

type handler struct {
	product  *handlers.Product
	category *handlers.Category
}

func New(cfg *Configuration) (*Server, error) {

	srv := &Server{
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

	//wrap with metrics
	crepo = repository.NewMetricsCategoryRepository(&metrics.AppMetrics.Category.Repository, crepo)
	repo = repository.NewMetricsProductRepository(&metrics.AppMetrics.Product.Repository, repo)

	svc := service.NewServiceMetrics(metrics.AppMetrics, service.NewService(crepo, repo))

	srv.handlers = &handler{
		product:  handlers.NewProduct(srv.errors, svc),
		category: handlers.NewCategory(srv.errors, svc),
	}
	if err := srv.registerRoutes(); err != nil {
		return nil, err
	}

	return srv, nil
}

func (s *Server) Serve(done <-chan interface{}, address string) error {
	s.server = &http.Server{Addr: address, Handler: s.router}
	s.info.Printf("Started serving at %s\n", address)
	go func(s *Server) {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			s.errors.Fatal(err)
		}
	}(s)

	pdone := make(chan interface{}, 1)
	if s.config.Prometheus.Push.Enabled && len(s.config.Prometheus.Push.Host) > 0 {
		pusher := push.New(s.config.Prometheus.Push.Host, s.config.Prometheus.Push.Job)
		go func(done <-chan interface{}, duration time.Duration, pusher *push.Pusher) {
			ticker := time.NewTicker(duration)
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					if err := pusher.Push(); err != nil {
						s.errors.Printf("pushing to prometheus: %v\n", err)
					}
				}
			}
		}(pdone, s.config.Prometheus.Push.Interval, pusher)
	}

	<-done
	pdone <- true
	return nil
}

func (s *Server) Shutdown(done chan<- interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.server.SetKeepAlivesEnabled(false)
	if err := s.server.Shutdown(ctx); err != nil {
		s.errors.Printf("shutting down server: %v\n", err)
	}
	close(done)
}
