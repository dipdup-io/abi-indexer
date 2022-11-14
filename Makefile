-include .env
export $(shell sed 's/=.*//' .env)

indexer:
	cd cmd/indexer && go run . -c ../../build/dipdup.yml

build-proto:
	protoc -I=. --go-grpc_out=./pkg ./pkg/modules/grpc/proto/*.proto
	protoc -I=. --go_out=./pkg ./pkg/modules/grpc/proto/*.proto

build:
	docker-compose up -d -- build
