proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/usermngnt.proto

server:
	go run server/server.go

client:
	go run client/client.go

.PHONY: proto server client