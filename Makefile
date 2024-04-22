create-docket-network:
	docker network create bank-network
postgres:
	docker run --name postgres16 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine
createdb:
	docker exec -it postgres16  createdb --username=root --owner=root bank_app
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@postgres:5432/bank_app?sslmode=disable" -verbose up
migrateuplastest:
	migrate -path db/migration -database "postgresql://root:secret@postgres:5432/bank_app?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@postgres:5432/bank_app?sslmode=disable" -verbose down
migratedownlastest:
	migrate -path db/migration -database "postgresql://root:secret@postgres:5432/bank_app?sslmode=disable" -verbose down 1
dropdb:
	docker exec -it postgres16 psql -U root bank_app
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/daniel-vuky/golang-bank-app/db/sqlc Store
.PHONY: create-docket-network postgres createdb migrateup migrateuplastest migratedown migratedownlastest dropdb sqlc test server mock