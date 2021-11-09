gen:
	protoc --go_out=. \
    --go-grpc_out=. \
    proto/*.proto

clean:
	rm pb/*.go

server:
	go run cmd/server/main.go -port 8000

client:
	go run cmd/client/main.go -address 0.0.0.0:8000
