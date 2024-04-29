.PHONY: migrate protoc

migrate:
	go run ./cmd/migrator/main.go --storage_path=./storage/sso.db --migrations_path=./migrations --goose_driver=sqlite3 --migrations_table=.

protoc:
	protoc -I ./contract/proto./contract/proto/sso/sso.proto --go_out=./contract/gen/go --go_opt=paths=source_relative --go-grpc_out=./contract/gen/go/ --go-grpc_opt=paths=source_relative

auth:
	go run ./cmd/sso/main.go --config=./config/local.yaml