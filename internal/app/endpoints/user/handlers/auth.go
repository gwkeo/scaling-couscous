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

type AuthService interface {
	ReadByEmail(context.Context, string) (*models.User, error)
}

type AuthHandler struct {
	service AuthService
	config  config.Config
}

func NewAuthHandler(service AuthService, config config.Config) *AuthHandler {
	return &AuthHandler{
		service: service,
		config:  config,
	}
}

func (h *AuthHandler) Authenticate(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		var user models.User
		err = json.Unmarshal(data, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userFromDB, err := h.service.ReadByEmail(ctx, user.Email)
		if err != nil {
			http.Error(w, err.Error()+" "+user.Email, http.StatusNotFound)
			return
		}

		if userFromDB == nil || userFromDB.Password != user.Password {
			http.Error(w, "wrong password", http.StatusUnauthorized)
			return
		} else {
			token, err := utils.GenerateToken(h.config.Secret, userFromDB.ID, userFromDB.Role)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("x-auth-token", fmt.Sprintf("Bearer %s", token))
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}
