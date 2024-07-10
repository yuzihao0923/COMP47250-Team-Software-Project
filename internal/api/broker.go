package api

import (
	"COMP47250-Team-Software-Project/internal/auth"
	"COMP47250-Team-Software-Project/internal/database"
	"COMP47250-Team-Software-Project/pkg/pool"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"net/http"

	"COMP47250-Team-Software-Project/internal/redis"
)

type HandlerResult struct {
	Data  interface{}
	Error error
}

func RegisterHandlers(mux *http.ServeMux, workerPool *pool.WorkerPool, db *database.MongoDB, rsi *redis.RedisServiceInfo) {
	jwtMiddleware := auth.JWTAuthMiddleware

	mux.Handle("/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultChan := make(chan HandlerResult)
		defer close(resultChan)

		workerPool.Submit(pool.JobFunc(func() {
			resultChan <- HandleLogin(db, r)
		}))

		// Handle result from worker
		result := <-resultChan
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusUnauthorized)
		} else {
			w.Header().Set("Content-Type", "application/json")
			serializer.JSONSerializerInstance.SerializeToWriter(result.Data, w)
		}
	}))

	mux.Handle("/produce", jwtMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultChan := make(chan HandlerResult)
		defer close(resultChan)

		workerPool.Submit(pool.JobFunc(func() {
			resultChan <- HandleProduce(rsi, r)
			// 模拟 broker 出现错误
			// resultChan <- HandlerResult{
			// 	Data:  nil,
			// 	Error: fmt.Errorf("broker cant handle produce"),
			// }
		}))

		result := <-resultChan
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			serializer.JSONSerializerInstance.SerializeToWriter(result.Data, w)
		}
	})))

	mux.Handle("/register", jwtMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultChan := make(chan HandlerResult)
		defer close(resultChan)

		workerPool.Submit(pool.JobFunc(func() {
			resultChan <- HandleRegister(rsi, r)
		}))

		result := <-resultChan
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			serializer.JSONSerializerInstance.SerializeToWriter(result.Data, w)
		}
	})))

	mux.Handle("/consume", jwtMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultChan := make(chan HandlerResult)
		defer close(resultChan)

		workerPool.Submit(pool.JobFunc(func() {
			resultChan <- HandleConsume(rsi, r)
		}))

		result := <-resultChan
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			serializer.JSONSerializerInstance.SerializeToWriter(result.Data, w)
		}
	})))

	mux.Handle("/ack", jwtMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultChan := make(chan HandlerResult)
		defer close(resultChan)

		workerPool.Submit(pool.JobFunc(func() {
			resultChan <- HandleACK(rsi, r)
		}))

		result := <-resultChan
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			serializer.JSONSerializerInstance.SerializeToWriter(result.Data, w)
		}
	})))
	// jwtMiddleware := auth.JWTAuthMiddleware

	// mux.Handle("/login", http.HandlerFunc(HandleLogin))
	// mux.Handle("/produce", jwtMiddleware(http.HandlerFunc(HandleProduce)))
	// mux.Handle("/register", jwtMiddleware(http.HandlerFunc(HandleRegister)))
	// mux.Handle("/consume", jwtMiddleware(http.HandlerFunc(HandleConsume)))
	// mux.Handle("/ack", jwtMiddleware(http.HandlerFunc(HandleACK)))
}
