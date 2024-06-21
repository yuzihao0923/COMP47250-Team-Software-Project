package api

import (
	"COMP47250-Team-Software-Project/internal/auth"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"fmt"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	err := serializer.JSONSerializerInstance.DeserializeFromReader(r.Body, &creds)
	if err != nil {
		log.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	// User credentials for different roles
	users := map[string]string{
		"broker":   "123",
		"consumer": "123",
		"producer": "123",
	}

	expectedPassword, ok := users[creds.Username]
	if !ok || creds.Password != expectedPassword {
		log.WriteErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid username or password"))
		return
	}

	token, err := auth.GenerateJWT(creds.Username)
	if err != nil {
		log.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	serializer.JSONSerializerInstance.SerializeToWriter(map[string]string{
		"token": token,
	}, w)
}
