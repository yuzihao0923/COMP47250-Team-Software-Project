package api

import (
	"COMP47250-Team-Software-Project/internal/auth"
	"COMP47250-Team-Software-Project/internal/database"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"context"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleLogin(r *http.Request) HandlerResult {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := serializer.JSONSerializerInstance.DeserializeFromReader(r.Body, &creds)
	if err != nil {
		return HandlerResult{nil, errors.New("failed to parse login request")}
	}

	var user struct {
		Username string `bson:"username"`
		Password string `bson:"password"`
		Role     string `bson:"role"`
	}
	err = database.UsersCollection.FindOne(context.TODO(), bson.M{"username": creds.Username}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return HandlerResult{nil, errors.New("this username is not valid, please try again")}
		} else {
			return HandlerResult{nil, err}
		}
	}

	if user.Password != creds.Password {
		return HandlerResult{nil, errors.New("this password is incorrect, please try again")}
	}

	token, err := auth.GenerateJWT(creds.Username)
	if err != nil {
		return HandlerResult{nil, errors.New("failed to generate token")}
	}

	// w.Header().Set("Content-Type", "application/json")
	// serializer.JSONSerializerInstance.SerializeToWriter(map[string]string{
	// 	"token":    token,
	// 	"username": creds.Username,
	// 	"role":     user.Role,
	// }, w)
	data := map[string]string{
		"token":    token,
		"username": creds.Username,
		"role":     user.Role,
	}
	return HandlerResult{Data: data, Error: nil}
}
