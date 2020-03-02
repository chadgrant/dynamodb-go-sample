package server

import (
	"expvar"
	"net/http"

	"github.com/chadgrant/go-http-infra/infra"
	"github.com/chadgrant/go-http-infra/infra/schema"
	"github.com/rs/cors"
)

func (s *Server) registerRoutes() error {

	gorillaW := func(str string, w http.HandlerFunc) {
		s.router.HandleFunc(str, w)
	}

	if err := infra.RegisterInfraHandlers(gorillaW, s.healthChecks, s.registry); err != nil {
		return err
	}
	if err := s.registry.AddDirectory("./schema"); err != nil {
		return err
	}
	var err error
	s.validator, err = schema.NewValidator(s.registry)
	if err != nil {
		return err
	}

	v := s.validator
	r := s.router
	ph := s.handlers.product
	ch := s.handlers.category

	r.Handle("/debug/vars", expvar.Handler())

	r.HandleFunc("/categories",
		v.Produces("http://schemas.sentex.io/store/categories.json", ch.GetAll),
	).Methods(http.MethodGet)

	r.HandleFunc("/products/{category:[A-Za-z]+}",
		v.Produces("http://schemas.sentex.io/store/product.paged.json", ph.GetPaged),
	).Methods(http.MethodGet)

	s.router.HandleFunc("/products/",
		v.Consumes("http://schemas.sentex.io/store/product-base.json", ph.Add),
	).Methods(http.MethodPost)

	r.HandleFunc("/product/{id:[a-z0-9\\-]{36}}",
		v.Consumes("http://schemas.sentex.io/store/product-base.json", ph.Upsert),
	).Methods(http.MethodPut)

	r.HandleFunc("/product/{id:[a-z0-9\\-]{36}}",
		v.Produces("http://schemas.sentex.io/store/product.json", ph.Get),
	).Methods(http.MethodGet)

	r.HandleFunc("/product/{id:[a-z0-9\\-]{36}}", ph.Delete).Methods(http.MethodDelete)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./docs/swagger/")))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		ExposedHeaders:   []string{"Location"},
		MaxAge:           86400,
	})

	r.Use(c.Handler)

	return nil
}
