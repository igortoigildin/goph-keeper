PROJECT_DIR = $(shell pwd)
PROJECT_BIN:=$(PROJECT_DIR)/bin
$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
PATH := $(PROJECT_BIN):$(PATH)

GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint

install-linter:
	### INSTALL GOLANGCI_LINT ###
	[ -f $(PROJECT_BIN)/golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v1.62.2

lint:
	make install-linter
	### RUN GOLANGCI_LINT ###
	$(GOLANGCI_LINT) run ./... --config=.golangci.yml

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2 
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	GOBIN=$(PROJECT_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.24.0

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/pressly/goose/v3/cmd/goose@v3.24.0

generate:
	make generate-upload-api
	make generate-auth-api

generate-upload-api:
	mkdir -p pkg/upload_v1
	protoc --proto_path api/upload_v1 \
	--go_out=pkg/upload_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/upload_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/upload_v1/upload.proto

generate-auth-api:
	mkdir -p pkg/auth_v1
	protoc --proto_path api/auth_v1 \
	--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/auth_v1/auth.proto

# make generate-access-api:
# 	mkdir -p pkg/access_v1
# 	protoc --proto_path api/access_v1 \
# 	--go_out=pkg/access_v1 --go_opt=paths=source_relative \
# 	--plugin=protoc-gen-go=bin/protoc-gen-go \
# 	--go-grpc_out=pkg/access_v1 --go-grpc_opt=paths=source_relative \
# 	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
# 	api/access_v1/access.proto

# These are the default values for the test database. They can be overridden
PG_DATABASE_NAME ?= test-db
PG_PORT ?= 54321
PG_PASSWORD ?= test-password
PG_USER ?= test-user
GOOSE_DRIVER ?= postgres

### RUN Goose migrations ###
-include .env
create-migration:
	$(PROJECT_BIN)/goose -dir ${MIGRATION_DIR} create add_users_table sql
	$(PROJECT_BIN)/goose -dir ${MIGRATION_DIR} create add_orders_table sql

migration-status:
	$(PROJECT_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

migration-up:
	$(PROJECT_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

migration-down:
	$(PROJECT_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v