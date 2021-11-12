gen:
	protoc --proto_path=proto proto/*.proto  --go_out=. --go-grpc_out=. --grpc-gateway_out=. --openapiv2_out=:swagger

clean:
	rm pb/*.go

grpc:
	go run cmd/server/main.go -grpc-port 8000 -server-type grpc

rest:
	go run cmd/server/main.go -http-port 8001 -server-type rest

both:
	go run cmd/server/main.go -grpc-port 8000 -http-port 8001

client:
	go run cmd/client/main.go -address 0.0.0.0:8000
