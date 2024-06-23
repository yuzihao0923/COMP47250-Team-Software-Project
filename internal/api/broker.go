package api

import (
	"COMP47250-Team-Software-Project/internal/auth"
	"COMP47250-Team-Software-Project/pkg/pool"
	"net/http"
)

func RegisterHandlers(mux *http.ServeMux, workerPool *pool.WorkerPool) {
	jwtMiddleware := auth.JWTAuthMiddleware

	mux.Handle("/login", jwtMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		workerPool.Submit(pool.JobFunc(func() {
			HandleLogin(w, r)
		}))
	})))

	mux.Handle("/produce", jwtMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		workerPool.Submit(pool.JobFunc(func() {
			HandleProduce(w, r)
		}))
	})))

	mux.Handle("/register", jwtMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		workerPool.Submit(pool.JobFunc(func() {
			HandleRegister(w, r)
		}))
	})))

	mux.Handle("/consume", jwtMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		workerPool.Submit(pool.JobFunc(func() {
			HandleConsume(w, r)
		}))
	})))

	mux.Handle("/ack", jwtMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		workerPool.Submit(pool.JobFunc(func() {
			HandleACK(w, r)
		}))
	})))
}
