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
	gotestsum --junitfile report.xml --format testname -- ./...

.PHONY: before-commit
before-commit: test ## run all checks before commit
	sqlfluff fix -n --disable-progress-bar --dialect postgres migrations/*.sql
	sqlfluff lint -n --disable-progress-bar --dialect postgres migrations/*.sql
	@golangci-lint run --timeout=10m
	@echo "Using config file: ${CONFIG_FILE}"
	@CONFIG_FILE=${CONFIG_FILE} gotestsum --junitfile report.xml --format testname ./...