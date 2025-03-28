.PHONY: build run stop build_and_run run_tests run_migrations

build:
	docker compose build

run:
	docker compose up -d

stop:
	docker compose down

build_and_run: build run

run_tests:
	docker compose run --rm app_test

run_migrations:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="host=localhost port=5432 password=postgres user=postgres dbname=tasks sslmode=disable" goose -dir ./task-service/migrations up