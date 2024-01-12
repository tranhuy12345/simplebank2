postgres15:
	docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecret -d postgres:15-alpine
createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres15 dropdb simple_bank
migrationup:
	migrate -path db/migration -database "postgresql://root:mysecret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrationdown:
	migrate -path db/migration -database "postgresql://root:mysecret@localhost:5433/simple_bank?sslmode=disable" -verbose down
opencommand:
	docker exec -it postgres15 psql -U root -d simple_bank
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
.PHONY: createdb dropdb migrationup migrationdown