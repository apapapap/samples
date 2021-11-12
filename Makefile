gen:
	protoc --proto_path=proto proto/*.proto  --go_out=. --go-grpc_out=. --grpc-gateway_out=. --openapiv2_out=:swagger

clean:
	rm pb/*.go

server:
	go run cmd/server/main.go -port 8000

rest:
	go run cmd/server/main.go -port 8001 -server-type rest

client:
	go run cmd/client/main.go -address 0.0.0.0:8000
