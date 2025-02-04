package handlers

import (
	"context"
	"encoding/json"
	"github.com/gwkeo/scaling-couscous/internal/app/repo/models"
	"net/http"
	"strconv"
)

type ReadUserInterface interface {
	Read(context.Context, int64) (*models.User, error)
	ReadByEmail(context.Context, string) (*models.User, error)
}

type ReadUserHandler struct {
	context context.Context
	service ReadUserInterface
}

func NewReadUserHandler(context context.Context, service ReadUserInterface) *ReadUserHandler {
	return &ReadUserHandler{
		context: context,
		service: service,
	}
}

func (h *ReadUserHandler) Read() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := h.service.Read(h.context, idInt64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(response)
	}
}
