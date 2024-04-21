VERSION ?= $(shell git describe --abbrev=0 --tags)
POSTGRES_URL ?= "postgres://test:changeme123@localhost:5432/test?sslmode=disable"
REPONAME ?= fatimalkaus
IMAGE_NAME ?= stack-server
DOCKER_TAG ?= latest
LD_FLAGS ?=
.PHONY: migrate-up
migrate-up:
	@migrate -database $(POSTGRES_URL) -path migrations/ up

.PHONY: migrate-down
migrate-down:
	@migrate -database $(POSTGRES_URL) -path migrations/ down

.PHONY: build
build:
	@go build -o bin/stack-server -ldflags "-X main.version=$(VERSION) -linkmode external -extldflags '-static' -s -w ${LD_FLAGS}" ./cmd/server/main.go

.PHONY: docker-build
docker-build:
	docker build --no-cache --tag $(REPONAME)/$(IMAGE_NAME):$(DOCKER_TAG) .

.PHONY: install-tools
install-tools:
	go install --tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

.PHONY: lint
lint:
	golangci-lint run -v