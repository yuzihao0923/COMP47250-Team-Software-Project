package api

import (
	"COMP47250-Team-Software-Project/internal/auth"
	"net/http"
)

func RegisterHandlers(mux *http.ServeMux) {
	jwtMiddleware := auth.JWTAuthMiddleware

	mux.Handle("/login", http.HandlerFunc(HandleLogin))
	mux.Handle("/produce", jwtMiddleware(http.HandlerFunc(HandleProduce)))
	mux.Handle("/register", jwtMiddleware(http.HandlerFunc(HandleRegister)))
	mux.Handle("/consume", jwtMiddleware(http.HandlerFunc(HandleConsume)))
	mux.Handle("/ack", jwtMiddleware(http.HandlerFunc(HandleACK)))
}