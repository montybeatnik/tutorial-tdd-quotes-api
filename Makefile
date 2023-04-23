SHELL := /bin/bash

db-up:
	docker volume create quotes_db
	docker run --name quotes_db \
	-e POSTGRES_PASSWORD=password \
	-v quotes_db:/var/lib/postgresql/data \
	-p 5432:5432 \
	-d postgres

db-sleep:
	sleep 3

db-create-quotes_db:
	docker exec quotes_db psql -U postgres -c "create database quotes_db"

db-create-table:
	docker exec quotes_db psql -U postgres -d quotes_db -c "CREATE TABLE IF NOT EXISTS quotes(id SERIAL PRIMARY KEY, author TEXT UNIQUE, message TEXT UNIQUE);"

db-all-up: db-up db-sleep db-create-inventory db-create-table

db-shell:
	docker exec -it quotes_db psql -U postgres -d quotes_db

db-clear:
	docker exec quotes_db psql -U postgres -d quotes_db -c "DELETE FROM quotes;"

db-down:
	docker container rm -f quotes_db
