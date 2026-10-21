package main

import (
	"fmt"
	"lesson/26/configs"
	"lesson/26/internal/auth"
	"lesson/26/pkg/db"
	"net/http"
)

func main() {
	cfg := configs.LoadConfig()
	_ = db.NewDb(cfg)
	router := http.NewServeMux()
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: cfg,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server is listening on port 8081...")
	server.ListenAndServe()
}
