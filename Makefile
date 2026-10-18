include .env

export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d postgres

env-plg-up:
	@docker compose up -d loki promtail grafana

env-plg-down:
	@docker compose stop loki promtail grafana

api-up:
	@docker compose up -d --build api

env-api-down:
	@docker compose down api

#Stop local env
env-down:
	@docker compose down

#Build local env
env-build:
	@docker compose build postgres

env-cleanup:
	@docker compose down postgres && \
	rm -rf postgres_data

migrate-action:
	@docker compose run --rm migrations \
		-path /migrations \
		-database postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable "$(action)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-redo:
	@make migrate-down && make migrate-up

migrate-force:
	@make migrate-action action="force $(version)"

migrate-local-action:
	@migrate -path ./migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" "$(action)"

migrate-local-up:
	@make migrate-local-action action=up

migrate-local-down:
	@make migrate-local-action action=down

env-port-forward:
	@docker compose up -d port-forwarder

env-port-down:
	@docker compose down port-forwarder

swagger-gen:
	@go run github.com/swaggo/swag/cmd/swag@v1.16.3 init -g cmd/api/main.go -o docs

bruno-sync:
	@echo "Syncing Bruno collection..."
	@bru import openapi --source docs/swagger.json --output ./bruno_collection --collection-name "Booking API"
	@echo "Done!"

kill-port:
	@lsof -t -i :$(or $(port),8080) | xargs kill -9 2>/dev/null || true

# docker compose logs -f api 
# вывод логов бэкенда