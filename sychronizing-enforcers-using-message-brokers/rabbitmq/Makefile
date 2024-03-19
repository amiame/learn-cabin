all:
	@docker compose up -d

rabbitmq:
	@docker compose up -d rabbitmq

producer:
	@go run ./producer

consumer:
	@go run ./consumer

clean:
	@docker compose down
	-pkill -f "go run ./producer" > /dev/null
	-pkill -f "go run ./consumer" > /dev/null

.PHONY: producer consumer
