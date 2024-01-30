postgres15:
	docker run --name postgres15.1 --network bank_network -p 5434:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecret -d postgres:15-alpine
createdb:
	docker exec -it postgres15.1 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres15.1 dropdb simple_bank
migrationup:
	migrate -path db/migration -database "postgres://root:mysecret@postgres:5432/simple_bank?sslmode=disable" -verbose up
migrationup1:
	migrate -path db/migration -database "postgresql://root:mysecret@localhost:5434/simple_bank?sslmode=disable" -verbose up 1
migrationdown:
	migrate -path db/migration -database "postgresql://root:mysecret@localhost:5434/simple_bank?sslmode=disable" -verbose down
migrationdown1:
	migrate -path db/migration -database "postgresql://root:mysecret@localhost:5434/simple_bank?sslmode=disable" -verbose down 1
opencommand:
	docker exec -it postgres15.1 psql -U root -d simple_bank
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
createmagration:
	migrate create -ext sql -dir db/migration -seq add_users
.PHONY: createdb dropdb migrationup migrationdown tst server sqlc migrationdown1 migrationup1