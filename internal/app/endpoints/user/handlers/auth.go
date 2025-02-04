package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gwkeo/scaling-couscous/internal/app/repo/models"
	"github.com/gwkeo/scaling-couscous/internal/config"
	"github.com/gwkeo/scaling-couscous/utils"
	"io"
	"net/http"
)

type AuthInterface interface {
	ReadByEmail(context.Context, string) (*models.User, error)
}

type AuthHandler struct {
	service AuthInterface
	context context.Context
	config  config.Config
}

func NewAuthHandler(service AuthInterface, context context.Context, config config.Config) *AuthHandler {
	return &AuthHandler{
		service: service,
		context: context,
		config:  config,
	}
}

func (h *AuthHandler) Authenticate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		defer r.Body.Close()

		var user models.User
		err = json.Unmarshal(data, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		userFromDB, err := h.service.ReadByEmail(r.Context(), user.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		if userFromDB == nil || userFromDB.Password != user.Password {
			http.Error(w, "wrong password", http.StatusUnauthorized)
		} else {
			token, err := utils.GenerateToken(h.config.Secret, userFromDB.ID, utils.User)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			w.Header().Set("x-auth-token", fmt.Sprintf("Bearer %s", token))
			w.WriteHeader(http.StatusOK)
		}
	}
}
