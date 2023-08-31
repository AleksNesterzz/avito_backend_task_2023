package getClientSeg

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	Status string   `json:"status"`
	Id     int      `json:"id,omitempty"`
	Seg    []string `json:"segments"`
	Error  string   `json:"error,omitempty"`
}

type Request struct {
	Id int `json:"id"`
}

type SegGetter interface {
	GetClientSeg(id int) ([]string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=SegGetter

func New(seggetter SegGetter) http.HandlerFunc {
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
		strings, err := seggetter.GetClientSeg(req.Id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			render.JSON(w, r, Response{
				Status: "error",
				Error:  "failed to get client segments"})

			return
		}

		render.JSON(w, r, Response{
			Status: "OK",
			Id:     req.Id,
			Seg:    strings,
		})
	}
}
