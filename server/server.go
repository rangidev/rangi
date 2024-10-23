package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/schema"

	"github.com/rangidev/rangi/admin"
	"github.com/rangidev/rangi/blueprint"
	"github.com/rangidev/rangi/config"
)

type Server struct {
	config            *config.Config
	server            *http.Server
	schemaDecoder     *schema.Decoder
	adminTemplates    *admin.Templates
	adminStaticServer http.Handler
	collectionLoader  *blueprint.CollectionLoader
}

func New(config *config.Config) (*Server, error) {
	// Request handling
	schemaDecoder := schema.NewDecoder()
	// Admin
	adminTemplates, err := admin.NewTemplates(config)
	if err != nil {
		return nil, fmt.Errorf("could not create admin templates: %v", err)
	}
	adminStaticServer, err := admin.NewStaticServer(config)
	if err != nil {
		return nil, fmt.Errorf("could not create admin static server: %v", err)
	}
	// Collections
	collectionLoader := blueprint.NewCollectionLoader(config.BlueprintsPath)
	collections, err := collectionLoader.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not get collections: %v", err)
	}
	// Create tables
	err = config.DatabaseInstance.CreateTables(collections, collectionLoader)
	if err != nil {
		return nil, fmt.Errorf("could not create tables: %v", err)
	}
	return &Server{
		config:            config,
		schemaDecoder:     schemaDecoder,
		adminTemplates:    adminTemplates,
		adminStaticServer: adminStaticServer,
		collectionLoader:  collectionLoader,
	}, nil
}

func (s *Server) Start() error {
	// Create router
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	// Add Admin
	// Static admin files without access check
	router.Get("/admin/static/*", s.GetAdminStatic)
	router.Mount("/admin", createAdminRouter(s))
	// Create server
	s.server = &http.Server{
		Addr:    s.config.HostAndPort,
		Handler: router,
	}
	// Start server
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
