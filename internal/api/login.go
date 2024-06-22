package api

import (
	"COMP47250-Team-Software-Project/internal/auth"
	"COMP47250-Team-Software-Project/internal/database"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	var user struct {
		Username string `bson:"username"`
		Password string `bson:"password"`
		Role     string `bson:"role"`
	}
	err = database.UsersCollection.FindOne(context.TODO(), bson.M{"username": creds.Username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "this username is not valid, please try again", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if user.Password != creds.Password {
		http.Error(w, "this password is incorrect, please try again", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(creds.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	serializer.JSONSerializerInstance.SerializeToWriter(map[string]string{
		"token":    token,
		"username": creds.Username,
		"role":     user.Role,
	}, w)
}
