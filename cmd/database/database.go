package main

import (
	"COMP47250-Team-Software-Project/internal/database"
	"COMP47250-Team-Software-Project/internal/log"
	"context"
	"fmt"
)

func main() {
	_, err := database.NewMongoDB("mongodb://root:****@dds-bp1ca790bb54aae43.mongodb.rds.aliyuncs.com:3717,dds-bp1ca790bb54aae41.mongodb.rds.aliyuncs.com:3717", "comp47250", "users")
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
