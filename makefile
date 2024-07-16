.PHONY: start stop proxy broker redis web

start: redis proxy broker web

redis:
	@echo "Starting Redis servers..."
	@cd internal/redis-cluster && \
	for port in 6381 6382 6383 6384 6385 6386; do \
		redis-server redis-$$port.conf & \
	done
	@sleep 1

proxy:
	@echo "Starting proxy..."
	@cd cmd/proxyServer && go run proxy.go &

broker:
	@echo "Starting broker..."
	@cd cmd/broker && go run broker.go &
	@sleep 1

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
	@for port in 6381 6382 6383 6384 6385 6386; do \
		redis-cli -p $$port shutdown; \
	done
	$(MAKE) kill-broker &
	$(MAKE) kill-proxy
