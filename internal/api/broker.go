package api

import (
	"COMP47250-Team-Software-Project/internal/auth"
	"COMP47250-Team-Software-Project/pkg/pool"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"net/http"
)

type HandlerResult struct {
	Data  interface{}
	Error error
}

func RegisterHandlers(mux *http.ServeMux, workerPool *pool.WorkerPool) {
	jwtMiddleware := auth.JWTAuthMiddleware

	mux.Handle("/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultChan := make(chan HandlerResult)
		defer close(resultChan)

		workerPool.Submit(pool.JobFunc(func() {
			resultChan <- HandleLogin(r)
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
			resultChan <- HandleProduce(r)
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
			resultChan <- HandleRegister(r)
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
			resultChan <- HandleConsume(r)
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
			resultChan <- HandleACK(r)
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
