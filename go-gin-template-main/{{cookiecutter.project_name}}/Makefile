include .env

default: help

.PHONY: help
help: ## Show help
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

.PHONY: swag
swag: ## Generate swagger documentation
	swag init

.PHONY: migrate-up
migrate-up: ## Apply migrations up to last version
	goose -table ${DATABASE_SCHEMA}.goose_db_version -dir ./migrations postgres ${DATABASE_DSN} up

.PHONY: migrate-zero
migrate-zero: ## Revert migrations to zero version
	goose -table ${DATABASE_SCHEMA}.goose_db_version -dir ./migrations postgres ${DATABASE_DSN} down
