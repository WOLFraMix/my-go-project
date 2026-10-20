package main

import (
	"fmt"
	"lesson/25/configs"
	"lesson/25/internal/auth"
	"net/http"
)

func main() {
	cfg := configs.LoadConfig()
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
