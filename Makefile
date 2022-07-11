run: ### Run docker-compose
	docker-compose up --build -d app && docker-compose logs -f
.PHONY: run

down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: down

unit-test: ### Run unit test
	go test -v -cover -race ./internal/...
.PHONY: unit-test

integration-test: down ### Run docker-compose with integration test
	docker-compose --profile integration-test up --build --abort-on-container-exit --exit-code-from integration
.PHONY: integration-test

swag-v1: ### Init swag
	swag init -g internal/controller/http/v1/router.go
.PHONY: swag-v1

swag-run: swag-v1 ### Run swag
	go mod tidy && go mod download && \
	DISABLE_SWAGGER_HTTP_HANDLER='' GIN_MODE=debug CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run-swag

