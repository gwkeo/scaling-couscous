package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	appErrors "github.com/gwkeo/scaling-couscous/internal/app/errors"
	"github.com/gwkeo/scaling-couscous/internal/app/repo/models"
	"io"
	"net/http"
)

type CreateUserInterface interface {
	Create(context.Context, *models.User) (int64, error)
}

type CreateUserHandler struct {
	service CreateUserInterface
	context context.Context
}

func NewCreateUserHandler(context context.Context, service CreateUserInterface) CreateUserHandler {
	return CreateUserHandler{
		service: service,
		context: context,
	}
}

func (c CreateUserHandler) CreateUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user := &models.User{}
		err = json.Unmarshal(body, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := c.service.Create(c.context, user)
		if err != nil {
			if errors.Is(err, appErrors.ErrFailedToInsertUser) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error() + ": user already exists"))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("User created: %v", id)
		_, err = w.Write([]byte(response))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
