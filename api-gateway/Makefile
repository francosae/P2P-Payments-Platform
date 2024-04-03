# Variables
BINARY_NAME=api-gateway
DOCKER_IMAGE_NAME=api-gateway
DOCKER_TAG=latest

# Regenerate proto files
proto: 
	protoc pkg/**/pb/*.proto --go_out=. --go-grpc_out=.

#Run the server
server:
	go run cmd/main.go

# Generate swagger docs
swagger:
	swag init -g cmd/main.go --parseDependency --parseInternal
	
# Build the Docker image
docker-build:
	docker build -t ${DOCKER_IMAGE_NAME}:${DOCKER_TAG} .

# Run the Docker container (currently mounting the .env file at runtime)
docker-run:
	docker run -p 3000:3000 -v `pwd`/pkg/config/envs:/root/pkg/config/envs ${DOCKER_IMAGE_NAME}:${DOCKER_TAG}

.PHONY: build run docker-build docker-run docker-push clean
