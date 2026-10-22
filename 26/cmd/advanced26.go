package main

import (
	"fmt"
	"lesson/26/configs"
	"lesson/26/internal/auth"
	"lesson/26/internal/link"
	"lesson/26/pkg/db"
	"net/http"
)

func main() {
	cfg := configs.LoadConfig()
	db := db.NewDb(cfg)
	router := http.NewServeMux()

	linkRepository := link.NewLinkRepository(db)

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: cfg,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server is listening on port 8081...")
	server.ListenAndServe()
}
