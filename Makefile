build:
	@go mod download && go build -o bin/friendly cmd/main.go

run: build
	@./bin/friendly

migrate-create:
	@migrate create -ext sql -dir ./db/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@migrate -path ./db/migrations -database "postgres://postgres:postgres@127.0.0.1:5432/friendly_db?sslmode=disable" up

migrate-down:
	@migrate -path ./db/migrations -database "postgres://postgres:postgres@127.0.0.1:5432/friendly_db?sslmode=disable" down