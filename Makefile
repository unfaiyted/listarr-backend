.PHONY: swag build run docker-build docker-run test

# Variables
DOCKER_IMAGE := listarr-backend 
GO_VERSION := 1.23
ALPINE_VERSION := 3.19

# Local development commands
swag:
	swag init

build: swag
	CGO_ENABLED=0 go build -o main main.go

run: swag
	go run main.go

test:
	go test ./... -v

# Docker commands
docker-build: swag
	docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build
	docker run -p 8080:8080 $(DOCKER_IMAGE)
