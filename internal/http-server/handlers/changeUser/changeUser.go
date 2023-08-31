package changeUser

import (
	"avito_backend/internal/storage"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

type UserChanger interface {
	ChangeUser(addSeg []string, delSeg []string, id int) (string, error)
}

type Request struct {
	AddSeg []string `json:"addseg,omitempty"`
	DelSeg []string `json:"delseg,omitempty"`
	Id     int      `json:"id"`
}

type Response struct {
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UserChanger

func New(userchanger UserChanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req Request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "failed to decode request body"})
			return
		}

		_, err = userchanger.ChangeUser(req.AddSeg, req.DelSeg, req.Id)
		if errors.Is(err, storage.ErrSegExists) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "segment already exists"})

			return
		}
		if errors.Is(err, storage.ErrSegNotExists) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "segment not exists"})

			return
		}
		if errors.Is(err, storage.ErrUserNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "user not in segment"})

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
			Message: "segments were successfully changed",
		})
	}
}
