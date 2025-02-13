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

type CreateUserHandler struct{ service CreateUserInterface }

func NewCreateUserHandler(service CreateUserInterface) CreateUserHandler {
	return CreateUserHandler{service: service}
}

func (c CreateUserHandler) CreateUser(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
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

		user.Role = models.UserRole

		id, err := c.service.Create(ctx, user)
		if err != nil {
			if errors.Is(err, appErrors.ErrFailedToInsertUser) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("User created: %v", id)
		http.Error(w, response, http.StatusCreated)
	}
}
