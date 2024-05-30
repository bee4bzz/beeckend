MIGRATE := migrate -path=migrations/ -database "$(DATABASE_URL)"

.PHONY: default
default: help

.PHONY: devenv-start
devenv-start: ## start all dependencies to test and run the web api (DB, mail, ...)
	docker-compose -f docker-compose.yml up -d

.PHONY: devenv-stop
devenv-stop: ## stop all dependencies to test and run the web api (DB, mail, ...)
	docker-compose down --remove-orphans

.PHONY: test
test: ## run all tests
	gotestsum --junitfile report.xml --format testname -- -coverprofile=cover.out ./... 

.PHONY: before-commit
before-commit: test ## run all checks before commit
	sqlfluff fix -n --disable-progress-bar --dialect postgres migrations/*.sql
	sqlfluff lint -n --disable-progress-bar --dialect postgres migrations/*.sql
	@golangci-lint run --timeout=10m
	@echo "Using config file: ${CONFIG_FILE}"
	@CONFIG_FILE=${CONFIG_FILE} gotestsum --junitfile report.xml --format testname ./...

.PHONY: migrate-reset
migrate-reset: ## reset database and re-run all migrations
	@echo "Resetting database..."
	@$(MIGRATE) drop
	@echo "Running all database migrations..."
	@$(MIGRATE) up

.PHONY: swag docs api
docs: swag ## generate OpenAPI/Swagger specs
api: swag ## generate OpenAPI/Swagger specs
swag: ## generate OpenAPI/Swagger specs
	swag init -g cmd/rest-api-server/main.go --parseInternal --parseDependency --output api

.PHONY: swag-fmt
swag-fmt: ## format swag comments
	swag fmt -g cmd/rest-api-server/main.go