package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"lesson/26/configs"
	"lesson/26/internal/link"
)

func main() {
	cfg := configs.LoadConfig()
	db, err := gorm.Open(postgres.Open(cfg.Db.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&link.Link{}); err != nil {
		log.Fatal(err)
	}
	log.Println("migrations applied")
}
