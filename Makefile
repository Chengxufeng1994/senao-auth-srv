redis:
	docker run --name redis -d -p 6379:6379 --network=senao-network --restart=always redis:7-alpine

server:
	go run main.go

docs:
	swag init .

docker-build:
	docker build --no-cache -t senao-auth-srv:latest -f Dockerfile .

docker-start:
	docker run --rm -it --name senao-auth-srv -p 8000:8000 senao-auth-srv:latest

.PHONY: network redis docs server docker-build
