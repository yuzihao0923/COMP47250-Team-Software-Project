# Makefile for front-end test

.PHONY: start stop broker web redis

start: broker web redis

broker:
	@echo "Starting broker..."
	@go run ./cmd/broker/broker.go &

web:
	@echo "Starting web server..."
	@cd web-app && npm start &

redis:
	@echo "Starting Redis servers..."
	@for conf_file in internal/redis-cluster/redis-6381.conf internal/redis-cluster/redis-6382.conf internal/redis-cluster/redis-6383.conf internal/redis-cluster/redis-6384.conf internal/redis-cluster/redis-6385.conf internal/redis-cluster/redis-6386.conf; do \
		redis-server $$conf_file & \
		sleep 1; \
	done

stop:
	@echo "Stopping all services..."
	@pkill -f 'go run ./cmd/broker/broker.go'
	@pkill -f 'npm start'
	@pkill redis-server
