proto:
	mkdir -p pb
	protoc --plugin=protoc-gen-go=/mnt/c/Users/ap\ comp\ yk/go/bin/protoc-gen-go.exe \
		--plugin=protoc-gen-go-grpc=/mnt/c/Users/ap\ comp\ yk/go/bin/protoc-gen-go-grpc.exe \
		--proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		proto/*.proto

proto-image-service:
	mkdir -p pb
	protoc --plugin=protoc-gen-go=/mnt/c/Users/ap\ comp\ yk/go/bin/protoc-gen-go.exe \
		--plugin=protoc-gen-go-grpc=/mnt/c/Users/ap\ comp\ yk/go/bin/protoc-gen-go-grpc.exe \
		--proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		proto/image_service.proto

# Build gRPC server
build-grpc:
	go build -o grpc_server.exe grpc_server.go

# Run gRPC server
run-grpc: build-grpc
	./grpc_server.exe

# Build HTTP server (existing)
build-http:
	go build -o auth_service.exe main.go

# Run HTTP server (existing)
run-http: build-http
	./auth_service.exe

# Run both servers (in background)
run-all:
	@echo "Starting HTTP server on port 8080..."
	./auth_service.exe &
	@echo "Starting gRPC server on port 50051..."
	./grpc_server.exe

# Clean build files
clean:
	rm -f *.exe
	rm -rf pb/*

.PHONY: proto build-grpc run-grpc build-http run-http run-all clean