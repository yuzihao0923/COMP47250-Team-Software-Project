# Makefile for front-end test

.PHONY: start stop proxy broker redis create-cluster web kill-proxy kill-broker

start: proxy broker redis create-cluster web kill-proxy kill-broker

proxy:
	@echo "Starting proxy..."
	@go run ./cmd/proxyServer/proxy.go &

broker:
	@echo "Starting broker..."
	@go run ./cmd/broker/broker.go &

redis:
	@echo "Starting Redis servers..."
	@for port in 6381 6382 6383 6384 6385 6386; do \
		redis-server ./internal/redis-cluster/redis-$$port.conf; \
	done

create-cluster:
	@echo "Creating Redis cluster..."
	@redis-cli --cluster create localhost:6381 localhost:6382 localhost:6383 localhost:6384 localhost:6385 localhost:6386 --cluster-replicas 1 --cluster-yes

# web:
#   @echo "Starting web server..."
#   @cd web-app && npm start &

kill-proxy:
	@echo "Killing all proxy processes..."
	@ps aux | grep '[p]roxy' | awk '{print $$2}' | xargs kill

kill-broker:
	@echo "Killing all broker processes..."
	@ps aux | grep '[b]roker' | awk '{print $$2}' | xargs kill

stop:
	@echo "Stopping all services..."
	@pkill redis-server
	$(MAKE) kill-broker
	$(MAKE) kill-proxy
	# @pkill -f 'go run ./cmd/proxyServer/proxy.go'
	# @pkill -f 'go run ./cmd/broker/broker.go'
	# @pkill -f 'npm start'