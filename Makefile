-include .env
export $(shell sed 's/=.*//' .env)

indexer:
	cd cmd/indexer && go run . -c ../../build/dipdup.yml

build-proto:
	protoc \
		-I=${GOPATH}/src \
		--go-grpc_out=${GOPATH}/src \
		--go_out=${GOPATH}/src \
		${GOPATH}/src/github.com/dipdup-net/abi-indexer/pkg/modules/grpc/proto/*.proto

build:
	docker-compose up -d --build
