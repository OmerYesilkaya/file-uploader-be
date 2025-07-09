run:
	@go run -tags "sqlite3" ./cmd

build:
	@go build -tags "sqlite3" -o bin/server ./cmd

create-db:
	@mkdir -p data
	@touch data/app.db

migrate-up:
	@../migrate/migrate -database "sqlite3://data/app.db" -path internal/db/migrations up

migrate-up-force:
	@../migrate/migrate -database "sqlite3://data/app.db" -path internal/db/migrations force 1

migrate-down:
	@./migrate/migrate -database "sqlite3://data/app.db" -path internal/db/migrations down

migration_create:
	@migrate create -ext sql -dir internal/db/migrations -seq $(name)
