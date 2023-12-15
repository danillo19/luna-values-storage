package main

import (
	"context"
	"luna-values-storage/internal/common"
	"luna-values-storage/internal/config"
	repo "luna-values-storage/internal/repository/mongodb"
	"luna-values-storage/internal/service"
	transport "luna-values-storage/internal/transport/http"
	"luna-values-storage/internal/transport/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Read()

	logger := common.GetLogger(false)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, disconnect, err := repo.NewDB(ctx, cfg.MongoURL, false)
	if err != nil {
		panic(err)
		return
	}

	defer disconnect()

	variableService := service.NewVariableService(db.VariableRepository)
	valueService := service.NewValueService(db.ValueRepository)

	s3Client, err := common.NewS3Client(logger)
	if err != nil {
		panic(err)
	}

	httpResolver := transport.NewResolver(variableService, valueService, s3Client, logger)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  middleware.CreateCorsAllowOriginFunc(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	//variables
	router.Get("/storage/vars/get", httpResolver.GetVariable)
	router.Post("/storage/vars/set", httpResolver.SetVariable)
	router.Post("/storage/vars/delete", httpResolver.DeleteVariable)

	//values
	router.Get("/storage/values/get", httpResolver.GetValue)
	router.Post("/storage/values/set", httpResolver.SetValue)
	router.Post("/storage/values/delete", httpResolver.DeleteValue)
	router.Post("/storage/files/upload", httpResolver.UploadFile)

	logger.Infof("listening on http://%s:%s", cfg.Host, cfg.Port)
	err = http.ListenAndServe(cfg.Host+":"+cfg.Port, router)
	if err != nil {
		panic(err)
	}

	err = logger.Sync()
	if err != nil {
		panic(err)
	}

}
