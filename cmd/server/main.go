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
	grpcPort := flag.Int("grpc-port", 0, "port to start the server on")
	httpPort := flag.Int("http-port", 0, "port to start the server on")
	serverType := flag.String("server-type", "", "type of server (grpc/rest)")
	flag.Parse()

	userServer := service.NewUserServer(service.NewInMemoryUserStore())
	roleServer := roleService.NewRoleServer(roleService.NewInMemoryRoleStore())

	var err error
	if *serverType == "grpc" {
		err = runGRPCServer(userServer, roleServer, *grpcPort)
	} else if *serverType == "rest" {
		err = runRESTServer(userServer, roleServer, *httpPort)
	} else {
		runBoth(userServer, roleServer, *grpcPort, *httpPort)
	}

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}

func runGRPCServer(userServer pb.UserServiceServer, roleServer pb.RoleServiceServer, grpcPort int) error {
	address := fmt.Sprintf("0.0.0.0:%d", grpcPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userServer)
	pb.RegisterRoleServiceServer(grpcServer, roleServer)

	log.Printf("starting gRPC server at %d", grpcPort)
	return grpcServer.Serve(listener)
}

func runRESTServer(userServer pb.UserServiceServer, roleServer pb.RoleServiceServer, httpPort int) error {
	address := fmt.Sprintf("0.0.0.0:%d", httpPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	mux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterUserServiceHandlerServer(ctx, mux, userServer)
	if err != nil {
		return err
	}

	err = pb.RegisterRoleServiceHandlerServer(ctx, mux, roleServer)
	if err != nil {
		return err
	}

	log.Printf("starting REST server at %d ", httpPort)
	return http.Serve(listener, mux)
}

func runBoth(userServer pb.UserServiceServer, roleServer pb.RoleServiceServer, grpcPort int, httpPort int) error {
	go runGRPCServer(userServer, roleServer, grpcPort)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcAddress := fmt.Sprintf("0.0.0.0:%d", grpcPort)
	err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return err
	}

	err = pb.RegisterRoleServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return err
	}

	httpAddress := fmt.Sprintf("0.0.0.0:%d", httpPort)
	listener, err := net.Listen("tcp", httpAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	log.Printf("starting REST proxy server at %d", httpPort)
	return http.Serve(listener, mux)
}
