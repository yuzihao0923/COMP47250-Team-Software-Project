package main

import (
	"COMP47250-Team-Software-Project/internal/database"
	"COMP47250-Team-Software-Project/internal/log"
	"context"
	"fmt"
)

func main() {
	_, err := database.NewMongoDB("mongodb://localhost:27017", "comp47250", "users")
	if err != nil {
		log.LogError("InitDB", err.Error())
		return
	}
	defer func() {
		ctx := context.Background()
		db := database.GetDBClient()
		if db != nil {
			if err := db.Close(ctx); err != nil {
				log.LogError("InitDB", "Failed to close MongoDB connection: "+err.Error())
			}
		}
	}()
	fmt.Println("Database initialized successfully")
}
