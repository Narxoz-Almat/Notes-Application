DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/bookstore?sslmode=disable
MIGRATE_DIR ?= migrations
NAME ?= change_name

.PHONY: run test migrate-up migrate-down migrate-force migrate-create

run:
	go run .

test:
	go test ./...

migrate-up:
	migrate -path $(MIGRATE_DIR) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATE_DIR) -database "$(DATABASE_URL)" down 1

migrate-force:
	migrate -path $(MIGRATE_DIR) -database "$(DATABASE_URL)" force $(VERSION)

migrate-create:
	migrate create -ext sql -dir $(MIGRATE_DIR) -seq $(NAME)