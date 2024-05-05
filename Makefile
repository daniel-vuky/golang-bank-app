create-docket-network:
	docker network create bank-network
postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine
redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine
createdb:
	docker exec -it postgres16  createdb --username=root --owner=root bank_app
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_app?sslmode=disable" -verbose up
migrateuplastest:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_app?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_app?sslmode=disable" -verbose down
migratedownlastest:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_app?sslmode=disable" -verbose down 1
dropdb:
	docker exec -it postgres16 psql -U root bank_app
sqlc:
	sqlc generate
test:
	go test -v -cover -short ./...
server:
	go run main.go
mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/daniel-vuky/golang-bank-app/db/sqlc Store
proto:
	rm -rf pb/*.go
	rm -rf doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
    --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=bank_app \
    proto/*.proto
evans:
	evans -r repl --host localhost --port 9090
.PHONY: create-docket-network postgres redis createdb migrateup migrateuplastest migratedown migratedownlastest dropdb sqlc test server mock proto evans