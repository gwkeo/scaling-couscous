package handlers

import (
	"context"
	ctx "github.com/gorilla/context"
	"github.com/gwkeo/scaling-couscous/internal/app/repo/models"
	"net/http"
	"strconv"
)

type DeleteUserInterface interface {
	Read(context.Context, int64) (*models.User, error)
	Delete(context.Context, int64) error
}

type DeleteUserHandler struct {
	context context.Context
	service DeleteUserInterface
}

func NewDeleteUserHandler(context context.Context, service DeleteUserInterface) *DeleteUserHandler {
	return &DeleteUserHandler{
		context: context,
		service: service,
	}
}

func (h *DeleteUserHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		roleRaw := ctx.Get(r, "role")
		if roleRaw == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		//if utils.Role(int64(roleRaw.(float64))) != utils.Admin {
		//
		//	w.WriteHeader(http.StatusForbidden)
		//	return
		//}

		id := r.URL.Query().Get("id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := h.service.Read(r.Context(), idInt)
		if err != nil {
			if user == nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.service.Delete(r.Context(), idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
