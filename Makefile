DB ?= postgres://tasks:tasks@localhost:5432/tasks?sslmode=disable

.PHONY: migrate-up migrate-down migrate-new

migrate-up:
	migrate -path ./migrations -database "$(DB)" up

migrate-down:
	migrate -path ./migrations -database "$(DB)" down 1

migrate-new:
	@read -p "name: " name; migrate create -ext sql -dir migrations -seq $$name