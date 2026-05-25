-include .env
export

export PROJECT_ROOT=$(shell pwd)


env-up:
	@docker compose up -d

env-down:
	@docker compose down

create-out:
	@mkdir -p out/benchmarks

deploy-run:
	docker compose up -d --build rwbtask-app

deploy-stop:
	@docker compose stop rwbtask-app

start-http-test:
	@make create-out && \
	hey -z 30s -c 100 "http://localhost:8080/toplist" \
       | tee out/benchmarks/bench_toplist_c100.txt

start-nats-test:
	@make create-out && \
	go run ./cmd/loadgen -n 3000000 -workers 32 -url nats://localhost:4222 \
       | tee out/benchmarks/bench_nats_mixed_long.txt

app-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/mainapp/main.go