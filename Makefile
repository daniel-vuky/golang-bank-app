postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine
createdb:
	docker exec -it postgres16  createdb --username=root --owner=root bank_app
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_app?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_app?sslmode=disable" -verbose down
dropdb:
	docker exec -it postgres16 psql -U root bank_app
sqlc:
	sqlc generate
test:
	go test -v -cover ./...

.PHONY: postgres createdb migrateup migratedown dropdb sqlc test