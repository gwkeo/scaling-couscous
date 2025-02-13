package handlers

import (
	"context"
	"encoding/json"
	gCtx "github.com/gorilla/context"
	"github.com/gwkeo/scaling-couscous/internal/app/repo/models"
	"io"
	"net/http"
)

type UpdateUserInterface interface {
	Update(context.Context, *models.User) error
}

type UpdateUserHandler struct {
	service UpdateUserInterface
}

func NewUpdateUserHandler(userService UpdateUserInterface) *UpdateUserHandler {
	return &UpdateUserHandler{service: userService}
}

func (u UpdateUserHandler) UpdateUser(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		rawId := gCtx.Get(r, "id")
		if rawId == nil {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}
		id := int64(rawId.(float64))

		if id == 0 {
			http.Error(w, "id parameter should be integer", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var user *models.User
		err = json.Unmarshal(body, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user.ID = id

		err = u.service.Update(ctx, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Successfully updated", http.StatusOK)
	}
}
