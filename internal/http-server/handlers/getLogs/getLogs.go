package getLogs

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/render"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Path   string `json:"path,omitempty"`
}

type Request struct {
	Id   int    `json:"id"`
	Time string `json:"time"`
}

type LogGetter interface {
	GetLogs(id int, year_month string) ([][]string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=LogGetter

func New(loggetter LogGetter) http.HandlerFunc {
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

		strings, err := loggetter.GetLogs(req.Id, req.Time)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			render.JSON(w, r, Response{
				Status: "error",
				Error:  "failed to get client segments"})

			return
		}
		id_string := strconv.Itoa(req.Id)
		path_to_file := "./csv_output/logs" + id_string + "-" + req.Time + ".csv"
		f2, err := os.OpenFile(path_to_file, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("error opening file")
			return
		}
		defer f2.Close()
		//fmt.Println(strings)
		csv_file := csv.NewWriter(f2)

		err = csv_file.WriteAll(strings)
		if err != nil {
			fmt.Println("error writing into file")
			return
		}

		fp, err := filepath.Abs(path_to_file)
		if err != nil {
			fmt.Println("error extracting path to file")
			return
		}

		render.JSON(w, r, Response{
			Status: "OK",
			Path:   fp,
		})
	}
}
