BINARY_NAME=pigeomail
DOMAIN=pigeomail.ddns.net

.DEFAULT_GOAL := help

PWD := $(shell pwd)

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

certs:
	openssl req -newkey rsa:2048 -sha256 -nodes -keyout .deploy/cert.key -x509 -days 365 -out .deploy/cert.pem -subj "/C=US/ST=New York/L=Brooklyn/O=Example Brooklyn Company/CN=${DOMAIN}"
	chmod 775 .deploy/cert.key
	chmod 775 .deploy/cert.pem

gogen_update:
	@printf "\033[33mUpdating rmq generator...\033[0m\n"
	@docker pull ghcr.io/requiemofthesouls/gogen:latest

gen_rmq:  ## run rmq generator
	@printf "\033[33mGenerating rmq code...\033[0m\n"
	@docker run --rm -it -v $(PWD):/app ghcr.io/requiemofthesouls/gogen:latest go generate ./proto/rmq
	@gofmt -s -w ./pb

gen_buf: ## run bufbuild generator
	@printf "\033[33mGenerating code with buf...\033[0m\n"
	@buf generate
	@go mod vendor

gen: gen_rmq gen_buf
	@go mod vendor

format: ## format code
	@printf "\033[36mFormatting code...\033[0m\n"
	@gofmt -s -w .

update: ## update dependencies
	@printf "\033[36mUpdating dependencies...\033[0m\n"
	@go get -u
	@go mod vendor

lint: ## run linter
	@printf "\033[36mRunning linter...\033[0m\n"
	@golangci-lint run

migration: ## create migration
	@printf "\033[36mCreating migration...\033[0m\n"
	@go run main.go migrate create $(name) -c deploy/local/config.dev.yaml -m deploy/migrations/

migrate_down: ## run migration down
	@printf "\033[36mRunning migrations down...\033[0m\n"
	@go run main.go migrate down -c deploy/local/config.dev.yaml -m deploy/migrations/

migrate_up: ## run migration up
	@printf "\033[36mRunning migrations up...\033[0m\n"
	@go run main.go migrate up -c deploy/local/config.dev.yaml -m deploy/migrations/

reset_db: migrate_down migrate_up ## reset database