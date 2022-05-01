export SERVER_HTTP_HOST=localhost
export SERVER_HTTP_PORT=8081

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run:
	cd ./cmd/ && \
	go run -tags dev .

build:
	rm -rf ./bin/ && \
	mkdir -p ./bin/ && \
	go build -ldflags="-s -w" -trimpath -o bin/main cmd/main.go

lint:
	golangci-lint run

protoc-gen: protoc-gen-user-details protoc-gen-user-details-mock protoc-gen-vehicle-service protoc-gen-vehicle-service-mock

protoc-gen-user-details:
	protoc \
	--proto_path=proto \
	--go_out=proto \
	--go_opt=paths=source_relative \
	--go-grpc_out=proto \
	--go-grpc_opt=paths=source_relative \
	user_details.proto

protoc-gen-user-details-mock:
	mockgen -source=proto/vehicle_service_grpc.pb.go -destination=proto/mock/vehicle_service.grpc_pb_mock.go -package=mock

protoc-gen-vehicle-service:
	protoc \
	--proto_path=proto \
	--go_out=proto \
	--go_opt=paths=source_relative \
	--go-grpc_out=proto \
	--go-grpc_opt=paths=source_relative \
	vehicle_service.proto

protoc-gen-vehicle-service-mock:
	mockgen -source=proto/vehicle_service_grpc.pb.go -destination=proto/mock/vehicle_service.grpc_pb_mock.go -package=mock

swagger: swagger-fmt
	swag init -g ./cmd/main.go -pd --parseDepth 2

swagger-fmt:
	swag fmt -g ./cmd/main.go

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
