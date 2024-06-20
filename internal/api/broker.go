package api

import (
	"COMP47250-Team-Software-Project/internal/auth"
	"net/http"
	"os"
)

func RegisterHandlers() {
	var jwtMiddleware func(http.Handler) http.Handler

	// Check if the environment is development
	if os.Getenv("DEVELOPMENT") == "true" {
		jwtMiddleware = auth.MockJWTAuthMiddleware
	} else {
		jwtMiddleware = auth.JWTAuthMiddleware
	}

	http.Handle("/login", CORSMiddleware(http.HandlerFunc(HandleLogin)))
	http.Handle("/produce", CORSMiddleware(jwtMiddleware(http.HandlerFunc(HandleProduce))))
	http.Handle("/register", CORSMiddleware(jwtMiddleware(http.HandlerFunc(HandleRegister))))
	http.Handle("/consume", CORSMiddleware(jwtMiddleware(http.HandlerFunc(HandleConsume))))
	http.Handle("/logs", CORSMiddleware(jwtMiddleware(http.HandlerFunc(HandleLogs))))
	http.Handle("/ack", CORSMiddleware(jwtMiddleware(http.HandlerFunc(HandleACK))))
}
