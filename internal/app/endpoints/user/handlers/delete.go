package handlers

import (
	"context"
	gorillaContext "github.com/gorilla/context"
	"github.com/gwkeo/scaling-couscous/internal/app/repo/models"
	"net/http"
	"strconv"
)

type DeleteUserInterface interface {
	Read(context.Context, int64) (*models.User, error)
	Delete(context.Context, int64) error
}

type DeleteUserHandler struct {
	service DeleteUserInterface
}

func NewDeleteUserHandler(service DeleteUserInterface) *DeleteUserHandler {
	return &DeleteUserHandler{service: service}
}

func (h *DeleteUserHandler) Delete(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		roleRaw := gorillaContext.Get(r, "role")
		if roleRaw == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if models.Role(int64(roleRaw.(float64))) != models.AdminRole {

			w.WriteHeader(http.StatusForbidden)
			return
		}

		id := r.URL.Query().Get("id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := h.service.Read(ctx, idInt)
		if err != nil {
			if user == nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.service.Delete(ctx, idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
