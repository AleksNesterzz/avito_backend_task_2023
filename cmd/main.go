package main

import (
	"avito_backend/internal/config"
	"avito_backend/internal/http-server/handlers/changeUser"
	"avito_backend/internal/http-server/handlers/createSeg"
	"avito_backend/internal/http-server/handlers/deleteSeg"
	"avito_backend/internal/http-server/handlers/getClientSeg"
	"avito_backend/internal/http-server/handlers/getLogs"
	"avito_backend/internal/storage/postgres"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {

	//загрузка конфига

	cfg := config.MustLoad()
	storagePath := "postgres://" + cfg.Db.Username + ":" + os.Getenv("DB_PASSWORD") + "@" + cfg.Db.Host + ":" + cfg.Db.Port + "/" + cfg.Db.Dbname + "?sslmode=disable"
	//инициализация БД
	storage, err := postgres.New(storagePath)
	if err != nil {
		log.Fatal("failed to init storage\n", err)
		os.Exit(1) // or 'return' (optional)
	}

	//инициализация роутера
	router := chi.NewRouter()

	//обработчики запросов

	router.Post("/seg/new", createSeg.New(storage))

	router.Patch("/user", changeUser.New(storage))

	router.Delete("/seg/del", deleteSeg.New(storage))

	router.Get("/user", getClientSeg.New(storage))

	router.Get("/log", getLogs.New(storage))

	//настройка сервера

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}
	//запуск сервера
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("failed to start server")
	}

}
