package user

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/gwkeo/scaling-couscous/internal/app/endpoints/auth"
	"github.com/gwkeo/scaling-couscous/internal/app/endpoints/user/handlers"
	"github.com/gwkeo/scaling-couscous/internal/app/repo/models"
	"github.com/gwkeo/scaling-couscous/internal/config"
	"net/http"
)

type Service interface {
	Create(context.Context, *models.User) (int64, error)
	Read(context.Context, int64) (*models.User, error)
	Update(context.Context, *models.User) error
	Delete(context.Context, int64) error
	ReadByEmail(context.Context, string) (*models.User, error)
}

type Server struct {
	router  *mux.Router
	service Service
	context context.Context
	config  *config.Config
}

func New(context context.Context, service Service, config *config.Config) *Server {
	return &Server{
		router:  mux.NewRouter(),
		context: context,
		service: service,
		config:  config,
	}
}

func (s *Server) registerRoutes() {
	createHandler := handlers.NewCreateUserHandler(s.context, s.service)
	s.router.HandleFunc(`/register`, createHandler.CreateUser()).Methods(http.MethodPost)
	authHandler := handlers.NewAuthHandler(s.service, s.context, *s.config)
	s.router.HandleFunc(`/login`, authHandler.Authenticate()).Methods(http.MethodPost)

	subRouter := s.router.PathPrefix("/user").Subrouter()
	subRouter.Use(auth.AuthenticationMW(s.config.Secret))
	readHandler := handlers.NewReadUserHandler(s.context, s.service)
	subRouter.HandleFunc(`/user`, readHandler.Read()).Methods(http.MethodGet)
	updateHandler := handlers.NewUpdateUserHandler(s.context, s.service)
	subRouter.HandleFunc(`/user`, updateHandler.UpdateUser()).Methods(http.MethodPut)
	deleteHandler := handlers.NewDeleteUserHandler(s.context, s.service)
	subRouter.HandleFunc(`/user`, deleteHandler.Delete()).Methods(http.MethodDelete)
}

func (s *Server) Start() error {

	s.registerRoutes()
	if err := http.ListenAndServe("localhost:8080", s.router); err != nil {
		return err
	}
	return nil
}
