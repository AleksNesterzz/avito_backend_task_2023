package deleteSeg

import (
	"avito_backend/internal/storage"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type SegDeleter interface {
	DeleteSeg(name string) (int64, error)
}

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
	Name    string `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=SegDeleter

func New(segdeleter SegDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req Request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			fmt.Println("failed to decode request body", err)
			return
		}

		name := req.Name

		_, err = segdeleter.DeleteSeg(name)
		if errors.Is(err, storage.ErrSegNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "segment not found"})
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "failed to delete segment"})

			return
		}

		render.JSON(w, r, Response{
			Status:  "OK",
			Name:    name,
			Message: "segment was successfully deleted",
		})
	}
}
