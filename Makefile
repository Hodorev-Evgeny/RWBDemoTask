-include .env
export

export PROJECT_ROOT=$(shell pwd)


env-up:
	@docker compose up -d

enc-down:
	@docker compose down

app-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/mainapp/main.go