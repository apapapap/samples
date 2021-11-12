package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"ashish/user-mgmt/pb"
	"ashish/user-mgmt/service"
	roleService "ashish/user-mgmt/service/role"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "port to start the server on")
	serverType := flag.String("server-type", "grpc", "type of server (grpc/rest)")
	flag.Parse()

	userServer := service.NewUserServer(service.NewInMemoryUserStore())
	roleServer := roleService.NewRoleServer(roleService.NewInMemoryRoleStore())

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	if *serverType == "grpc" {
		err = runGRPCServer(userServer, roleServer, listener)
	} else {
		err = runRESTServer(userServer, roleServer, listener)
	}

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}

func runGRPCServer(userServer pb.UserServiceServer, roleServer pb.RoleServiceServer, listener net.Listener) error {
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userServer)
	pb.RegisterRoleServiceServer(grpcServer, roleServer)

	return grpcServer.Serve(listener)
}

func runRESTServer(userServer pb.UserServiceServer, roleServer pb.RoleServiceServer, listener net.Listener) error {
	mux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterUserServiceHandlerServer(ctx, mux, userServer)
	if err != nil {
		return err
	}

	err = pb.RegisterRoleServiceHandlerServer(ctx, mux, roleServer)
	if err != nil {
		return err
	}

	log.Printf("starting REST server at %s ", listener.Addr().String())
	return http.Serve(listener, mux)

}
