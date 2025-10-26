proto:
	mkdir -p pb
	protoc --plugin=protoc-gen-go=/mnt/c/Users/ap\ comp\ yk/go/bin/protoc-gen-go.exe \
		--plugin=protoc-gen-go-grpc=/mnt/c/Users/ap\ comp\ yk/go/bin/protoc-gen-go-grpc.exe \
		--proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		proto/*.proto

.PHONY: proto