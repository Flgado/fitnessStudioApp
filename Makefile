.DEFAULT_GOAL := run-migration

docker-up:
	docker-compose up -d
	@echo "Waiting for PostgreSQL to be ready..."
	@docker-compose exec -T postgres sh -c 'while ! pg_isready -q -h localhost -p 5432 -U postgres; do sleep 1; done'
	@echo "PostgreSQL is ready!"

run-migration: docker-up
	@echo "Building migration binary..."
	go build -o cmd/sqlmigrations/main ./cmd/sqlmigrations/main.go
	@echo "Running migration..."
	./cmd/sqlmigrations/main
