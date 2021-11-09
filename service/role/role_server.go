package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ashish/user-mgmt/pb"
)

// RoleServer is the server that provides role services
type RoleServer struct {
	pb.UnimplementedRoleServiceServer
	roleStore RoleStore
}

// NewRoleServer returns a new RoleServer
func NewRoleServer(roleStore RoleStore) *RoleServer {
	return &RoleServer{
		roleStore: roleStore,
	}
}

func (server *RoleServer) AddRole(ctx context.Context, req *pb.AddRoleRequest) (*pb.AddRoleResponse, error) {
	role := req.GetRole()
	log.Printf("Received a create-role request with name: %s", role.Name)

	err := server.roleStore.Save(role)
	if err != nil {
		errCode := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			errCode = codes.AlreadyExists
			log.Printf("cannot save role with name: %s since it already exists", role.Name)
		}
		return nil, status.Errorf(errCode, "cannot save role to store: %v", err)
	}

	log.Printf("saved role with name: %s", role.Name)

	res := &pb.AddRoleResponse{
		Name: role.Name.String(),
	}
	return res, nil
}

func (server *RoleServer) FindRole(ctx context.Context, req *pb.FindRoleRequest) (*pb.FindRoleResponse, error) {
	found, err := server.roleStore.Find(req.Name)
	if err != nil {
		errMsg := fmt.Sprintf("cannot find role: %v", err)
		log.Print(errMsg)
		return nil, status.Error(codes.Internal, errMsg)
	}
	if found == nil {
		errMsg := fmt.Sprintf("role %s is not found", req.Name)
		log.Print(errMsg)
		return nil, status.Error(codes.NotFound, errMsg)
	}
	res := &pb.FindRoleResponse{
		Role: found,
	}
	return res, nil
}

func (server *RoleServer) FindRoles(ctx context.Context, req *pb.FindRolesRequest) (*pb.FindRolesResponse, error) {
	rolesMap, err := server.roleStore.FindAll()
	if err != nil {
		errMsg := fmt.Sprintf("cannot find roles: %v", err)
		log.Print(errMsg)
		return nil, status.Error(codes.Internal, errMsg)
	}

	if len(rolesMap) == 0 {
		errMsg := "no roles found"
		log.Print(errMsg)
		return nil, status.Error(codes.NotFound, errMsg)
	}

	var roles []*pb.Role
	for _, role := range rolesMap {
		roles = append(roles, role)
	}

	res := &pb.FindRolesResponse{
		Roles: roles,
	}
	return res, nil
}
