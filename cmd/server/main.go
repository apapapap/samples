package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"ashish/user-mgmt/pb"
	"ashish/user-mgmt/service"
	roleService "ashish/user-mgmt/service/role"

	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "port to start the server on")
	flag.Parse()

	log.Printf("Server started on port: %d", *port)

	userServer := service.NewUserServer(service.NewInMemoryUserStore())
	roleServer := roleService.NewRoleServer(roleService.NewInMemoryRoleStore())

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userServer)
	pb.RegisterRoleServiceServer(grpcServer, roleServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}
