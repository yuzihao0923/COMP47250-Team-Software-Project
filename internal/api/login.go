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
    }
    err := serializer.JSONSerializerInstance.DeserializeFromReader(r.Body, &creds)
    if err != nil {
        log.WriteErrorResponse(w, http.StatusBadRequest, err)
        return
    }

    // This is just for demonstration purposes. In real life, you should validate
    // the username and password against a database or another authentication service.
    if creds.Username == "user" && creds.Password == "password" {
        token, err := auth.GenerateJWT(creds.Username)
        if err != nil {
            log.WriteErrorResponse(w, http.StatusInternalServerError, err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        serializer.JSONSerializerInstance.SerializeToWriter(map[string]string{
            "token": token,
        }, w)
    } else {
        log.WriteErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid username or password"))
    }
}
