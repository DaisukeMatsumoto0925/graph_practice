DB_HOST=db
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=graphql-app-development
DB_CONN=mysql://${DB_USER}:${DB_PASSWORD}@tcp\(${DB_HOST}:${DB_PORT}\)/${DB_NAME}
curdir := `pwd`

.PHONY: run
run:
	docker-compose up --build -d

.PHONY: start
start:
	docker-compose exec app realize start --run

.PHONY: generate
generate:
	docker-compose exec app go generate ./...

.PHONY: migrate-create
migrate-create:
	docker-compose exec app migrate create -ext sql -dir migrations ${FILENAME}

.PHONY: migrate-up
migrate-up:
	docker-compose exec app migrate --source file://migrations --database ${DB_CONN} up

.PHONY: migrate-down
migrate-down:
	docker-compose exec app migrate --source file://migrations --database ${DB_CONN} down 1

schema-spy:
	docker run --net="host" \
	-v "$(curdir)/docs:/output"  \
	--platform=linux/amd64 schemaspy/schemaspy:latest -all  \
	-t mysql -host 127.0.0.1:${DB_PORT} -db ${DB_NAME} -u ${DB_USER} -p ${DB_PASSWORD} -connprops useSSL\\=false

