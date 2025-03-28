.PHONY: start_services stop_services build_services start_with_tests

start_services: build_services
	docker compose up

stop_services:
	docker compose down

build_services:
	docker compose build

start_with_tests:
	docker compose --profile test up

run_migrations:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="host=localhost port=5432 password=postgres user=postgres dbname=tasks sslmode=disable" goose -dir ./task-service/migrations up