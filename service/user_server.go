package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ashish/user-mgmt/pb"
)

// UserServer is the server that provides user services
type UserServer struct {
	pb.UnimplementedUserServiceServer
	userStore UserStore
}

// NewUserServer returns a new UserServer
func NewUserServer(userStore UserStore) *UserServer {
	return &UserServer{
		userStore: userStore,
	}
}

func (server *UserServer) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	user := req.GetUser()
	log.Printf("Received a create-user request with id: %s", user.Id)

	if len(user.Id) > 0 {
		_, err := uuid.Parse(user.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "user ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate new user ID: %v", err)
		}
		user.Id = id.String()
	}
	err := server.userStore.Save(user)
	if err != nil {
		errCode := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			errCode = codes.AlreadyExists
			log.Printf("cannot save user with id: %s since it already exists", user.Id)
		}
		return nil, status.Errorf(errCode, "cannot save user to store: %v", err)
	}

	log.Printf("saved user with id: %s", user.Id)

	res := &pb.AddUserResponse{
		Id: user.Id,
	}
	return res, nil
}

func (server *UserServer) FindUser(ctx context.Context, req *pb.FindUserRequest) (*pb.FindUserResponse, error) {
	found, err := server.userStore.Find(req.Id)
	if err != nil {
		errMsg := fmt.Sprintf("cannot find user: %v", err)
		log.Print(errMsg)
		return nil, status.Error(codes.Internal, errMsg)
	}
	if found == nil {
		errMsg := fmt.Sprintf("user %s is not found", req.Id)
		log.Print(errMsg)
		return nil, status.Error(codes.NotFound, errMsg)
	}
	res := &pb.FindUserResponse{
		User: found,
	}
	return res, nil
}

func (server *UserServer) FindUsers(ctx context.Context, req *pb.FindUsersRequest) (*pb.FindUsersResponse, error) {
	usersMap, err := server.userStore.FindAll()
	if err != nil {
		errMsg := fmt.Sprintf("cannot find users: %v", err)
		log.Print(errMsg)
		return nil, status.Error(codes.Internal, errMsg)
	}

	if len(usersMap) == 0 {
		errMsg := "no users found"
		log.Print(errMsg)
		return nil, status.Error(codes.NotFound, errMsg)
	}

	var users []*pb.User
	for _, user := range usersMap {
		users = append(users, user)
	}

	res := &pb.FindUsersResponse{
		Users: users,
	}
	return res, nil
}
