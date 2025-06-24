include .env

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATE=migrate -path ./migrations -database "$(DB_URL)"

create-migrations:
	migrate create -ext sql -dir ./migrations -seq init_schema

migrate-up:
	docker-compose up -d api-server
	docker-compose exec api-server bash -c "$(MIGRATE) up"

migrate-down:
	docker-compose up -d api-server
	docker-compose exec api-server bash -c "$(MIGRATE) down"

start:
	docker-compose up

stop:
	docker-compose down

.PHONY: deps
deps:
	go mod tidy

.PHONY: swagger-init
swagger-init:
	@echo "Generate swagger gui"
	swag init -g  cmd/main.go

.PHONY: create-migrations migrate-up migrate-down start stop run 