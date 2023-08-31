package createSeg

import (
	"avito_backend/internal/storage"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type SegCreator interface {
	CreateSeg(name string) (int64, error)
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

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=SegCreator

func New(segcreator SegCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req Request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			fmt.Println("failed to decode request body", err)
			return
		}

		name := req.Name

		//fmt.Println(id, name)

		_, err = segcreator.CreateSeg(name)
		if errors.Is(err, storage.ErrSegExists) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "segment already exists"})

			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "failed to add segment"})

			return
		}

		render.JSON(w, r, Response{
			Status:  "OK",
			Name:    name,
			Message: "segment was successfully added",
		})
	}
}
