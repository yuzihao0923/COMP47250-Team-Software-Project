.PHONY: start stop proxy broker redis initdb web kill-broker1#consumers

start: broker #redis initdb proxy broker web #consumers

redis:
	@echo "Starting Redis servers..."
	@cd internal/redis-cluster && \
	for port in 6381 6382 6383 6384 6385 6386; do \
		redis-server redis-$$port.conf & \
	done


	# @cd tests/four-master-cluster/cluster && \
	# for port in 6387 6388 6389 6390 6391 6392 6393 6394; do \
	# 	redis-server redis-$$port.conf & \
	# done

	# @cd tests/five-master-cluster/cluster && \
	# for port in 6395 6396 6397 6398 6399 6400 6401 6402 6403 6404; do \
	# 	redis-server redis-$$port.conf & \
	# done

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


web:
	@echo "Starting web..."
	@cd web-app && npm start &

broker: broker1 broker2 broker3 broker4 

# consumers: consumer1 consumer2 consumer3

# consumer1:
# 	@echo "Starting consumer 1..."
# 	@cd cmd/consumer && go run consumer.go -stream mystream1 -username c1 -password 123 &

# consumer2:
# 	@echo "Starting consumer 2..."
# 	@cd cmd/consumer && go run consumer.go -stream mystream2 -username c2 -password 123 &

# consumer3:
# 	@echo "Starting consumer 3..."
# 	@cd cmd/consumer && go run consumer.go -stream mystream3 -username c3 -password 123 &

kill-proxy:
	@echo "Killing all proxy processes..."
	@ps aux | grep 'exe/[p]roxy' | awk '{print $$2}' | xargs kill

kill-broker:
	@echo "Killing all broker processes..."
	@ps aux | grep 'exe/[b]roker' | awk '{print $$2}' | xargs kill

kill-web:
	@echo "Killing web process on port 3000..."
	@lsof -i :3000 | awk 'NR>1 {print $$2}' | xargs kill

kill-broker1:
	@echo "Killing broker1 on port 8081..."
	@lsof -t -i:8081 | xargs kill -9 || echo "No process found on port 8081"

# kill-consumers:
# 	@echo "Killing all consumer processes..."
# 	@ps aux | grep 'exe/[c]onsumer' | awk '{print $$2}' | xargs kill

stop: 
	@echo "Stopping all services..."
	@for port in 6381 6382 6383 6384 6385 6386; do \
		redis-cli -p $$port shutdown; \
	done
	$(MAKE) kill-broker
	$(MAKE) kill-proxy
	$(MAKE) kill-web
