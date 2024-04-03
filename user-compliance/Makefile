# Variables
BINARY_NAME=user-compliance-service
DOCKER_IMAGE_NAME=user-compliance-service
DOCKER_TAG=latest

# Regenerate proto files
proto: 
	protoc pkg/pb/*.proto --go_out=. --go-grpc_out=.

#Run the server
server:
	go run cmd/main.go

# Build the Docker image
docker-build:
	docker build -t ${DOCKER_IMAGE_NAME}:${DOCKER_TAG} .

# Run the Docker container (currently mounting the .env file at runtime)
docker-run:
	docker run -p 50051:50051 -v `pwd`/pkg/config/envs:/root/pkg/config/envs ${DOCKER_IMAGE_NAME}:${DOCKER_TAG}

.PHONY: build run docker-build docker-run docker-push clean
