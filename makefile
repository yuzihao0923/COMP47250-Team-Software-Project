# Makefile for front-end test

.PHONY: start stop broker web

start: broker web

broker:
	@echo "Starting broker..."
	@go run ./cmd/broker/broker.go &

web:
	@echo "Starting web server..."
	@cd web-app && npm start

stop:
	@echo "Stopping all services..."
	@pkill -f 'go run ./cmd/broker/broker.go'
	@pkill -f 'npm start'