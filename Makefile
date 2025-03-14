.PHONY: start_services stop_services build_services test

start_services: build_services
	docker compose up

stop_services:
	docker compose down

build_services:
	docker compose build

test:
	pytest -v 