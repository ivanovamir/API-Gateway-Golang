.SILENT:

include .env

run:
	docker-compose up --build -d

run_dev:
	docker compose -f dev.docker-compose.yml up --build -d

stop:
	docker-compose stop