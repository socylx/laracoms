build:
	protoc --go_out=. --go_opt=paths=source_relative --micro_out=. --micro_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/demo/demo.proto
	docker build -t laracom-demo-service .

run:
	docker run -d -p 9091:9091 -e MICRO_SERVER_ADDRESS=:9091 -e MICRO_REGISTRY=mdns laracom-demo-service
