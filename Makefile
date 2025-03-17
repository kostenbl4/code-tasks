.PHONY: start_services stop_services build_services start_with_tests

start_services: build_services
	docker compose up

stop_services:
	docker compose down

build_services:
	docker compose build

start_with_tests:
	docker compose --profile test up --abort-on-container-exit --exit-code-from app_test

	