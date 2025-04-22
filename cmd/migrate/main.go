package main

import (
	"fmt"
	"log"
	"socket/internal/config"
	"socket/internal/core/domain"
	"socket/internal/database"
)

func main() {
	config := config.Load()
	db, err := database.New(config.DB)
	if err != nil {
		log.Fatalf("Failed to init  DB err: %v", err)
	}
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	err = db.AutoMigrate(
		&domain.Sample{},
		&domain.User{},
		&domain.Message{},
		&domain.Chat{},
		&domain.MessageQueue{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	} else {
		fmt.Println("Migrate Successfully")
	}
}
