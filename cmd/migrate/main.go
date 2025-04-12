package main

import (
	"fmt"
	"socket/internal/config"
	"socket/internal/core/domain"
	"socket/internal/database"
)

func main() {
	config := config.Load()
	db, err := database.New(config.DB)
	if err != nil {
		fmt.Println("Error to init db")
	}
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	err = db.AutoMigrate(
		&domain.Sample{},
	)
	if err != nil {
		fmt.Println("Cannot auto migrate")
	} else {
		fmt.Println("Migrate Successfully")
	}
}
