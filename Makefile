include .env

# Docker
docker-up:
	sudo docker-compose up -d

docker-down:
	sudo docker-compose down

docker-ps:
	sudo docker-compose ps

docker-psql:
	sudo docker-compose exec postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}

docker-start:
	sudo docker-compose start

docker-stop:
	sudo docker-compose stop

docker-logs:
	sudo docker-compose logs

docker-volumes:
	sudo docker volume ls

volume-prune:
	sudo docker volume prune

# Migrate
migrate-init:
	migrate create -ext=sql -dir=${MIGRATIONS_PATH} -seq init

migrate-up:
	migrate -path=${MIGRATIONS_PATH} -database "${POSTGRES_PATH}?sslmode=disable" -verbose up

migrate-down:
	migrate -path=${MIGRATIONS_PATH} -database "${POSTGRES_PATH}?sslmode=disable" -verbose down	

# Go
run:
	go run main.go