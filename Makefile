pullpostgres:
	docker pull postgres:12-alpine

postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
	
dropdb:
	docker exec -it postgres12 dropdb --username=root simple_bank

dbup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

dbdown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

fulltest:
	go clean -testcache
	go test -v -coverpkg=./... -coverprofile=profile.cov ./...
	go test -v -coverpkg=./... -coverprofile=profile.cov ./...
	go tool cover -func profile.cov

server:
	go run main.go

.PHONY: createdb dropdb postgres migrateup migratdown sqlc test