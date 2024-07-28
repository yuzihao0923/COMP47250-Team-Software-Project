.PHONY: start stop proxy broker redis initdb web broker1 broker2 broker3 broker4

start: redis initdb proxy broker

redis:
	@echo "Starting Redis servers..."
	@cd internal/redis-cluster && \
	for port in 6381 6382 6383 6384 6385 6386; do \
		redis-server redis-$$port.conf & \
	done
	@sleep 1

initdb:
	@echo "Initializing database..."
	@cd cmd/database && go run database.go

proxy:
	@echo "Starting proxy..."
	@cd cmd/proxyServer && go run proxy.go &

broker1:
	@echo "Starting broker 1..."
	@cd cmd/broker && go run broker.go -id broker1 &

broker2:
	@echo "Starting broker 2..."
	@cd cmd/broker && go run broker.go -id broker2 &

broker3:
	@echo "Starting broker 3..."
	@cd cmd/broker && go run broker.go -id broker3 &

broker4:
	@echo "Starting broker 4..."
	@cd cmd/broker && go run broker.go -id broker4 &

# broker5:
# 	@echo "Starting broker 5..."
# 	@cd cmd/broker && go run broker.go -id broker5 &

# broker6:
# 	@echo "Starting broker 6..."
# 	@cd cmd/broker && go run broker.go -id broker6 &
broker: broker1 broker2 broker3 broker4

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
