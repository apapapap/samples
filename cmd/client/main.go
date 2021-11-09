package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"math/rand"
	"time"

	"ashish/user-mgmt/pb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrAlreadyExists = errors.New("resource already exists")
var ErrNotFound = errors.New("resource not found")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	serverAddress := flag.String("address", "", "address to start the server on")
	flag.Parse()

	log.Printf("Dial server: %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	userClient := pb.NewUserServiceClient(conn)
	roleClient := pb.NewRoleServiceClient(conn)

	log.Print("\n================ 1. Role Add ================")
	// roleName := randomStringFromSet("OWNER", "ADMIN", "CUSTOMER")
	rolePermissions := []pb.Permissions{pb.Permissions_READ, pb.Permissions_READWRITE}
	role, err := addRole(roleClient, pb.Name_OWNER, rolePermissions)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			log.Print("role already exists")
		} else {
			log.Print("cannot create role: ", err)
		}
	} else {
		log.Print("role created with name: ", role)
	}

	log.Print("\n================ 2. List all roles ================")
	roles, err := fetchAllRoles(roleClient)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Print("roles not found")
		} else {
			log.Print("cannot find roles: ", err)
		}
	} else {
		for index, role := range roles {
			log.Printf("	Role: %d\n", index+1)
			log.Printf("		ID: %s\n", role.Name)
			log.Print("		Permissions: ", role.Permissions)
		}
	}

	log.Print("\n================ 3. Role Fetch ================")
	fetchedRole, err := fetchRole(roleClient, role)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Print("role not found")
		} else {
			log.Print("cannot find role: ", err)
		}
	} else {
		log.Print("role found")
		log.Print("		Name: ", fetchedRole.Name)
		log.Print("		Permissions: ", fetchedRole.Permissions)
	}

	log.Print("\n================ 4. User Add ================")
	name := randomStringFromSet("Ashish", "Pratik", "Rajat")
	userAddress := randomStringFromSet("Pune", "Mumbai", "Nashik")
	anotherRole := &pb.Role{
		Name:        pb.Name_ADMIN,
		Permissions: []pb.Permissions(pb.Permissions_READWRITE.String()),
	}
	userRoles := []*pb.Role{fetchedRole, anotherRole}
	userID, err := addUser(userClient, name, userAddress, userRoles)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			log.Print("user already exists")
		} else {
			log.Print("cannot add user: ", err)
		}
	} else {
		log.Print("user added with id: ", userID)
	}

	log.Print("\n================ 5. List all users ================")
	users, err := fetchAllUsers(userClient)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Print("users not found")
		} else {
			log.Print("cannot find users: ", err)
		}
	} else {
		for index, user := range users {
			log.Printf("\tUser: %d\n", index+1)
			log.Printf("\t\tID: %s\n", user.Id)
			log.Print("\t\tName: ", user.Name)
			log.Print("\t\tAddress: ", user.Address)
			log.Print("\t\tRoles: \n")
			for index, role := range user.Roles {
				log.Printf("\t\t\tRole %d: \n", index)
				log.Print("\t\t\t\tRole-Name: ", role.Name)
				log.Print("\t\t\t\tRole-permissions: ", role.Permissions)
			}
		}
	}

	log.Print("\n================ 6. User Fetch ================")
	user, err := fetchUser(userClient, userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Print("user not found")
		} else {
			log.Print("cannot find user: ", err)
		}
	} else {
		log.Print("user found")
		log.Print("\t\tID: ", user.Id)
		log.Print("\t\tName: ", user.Name)
		log.Print("\t\tAddress: ", user.Address)
		log.Print("\t\tRoles: \n")
		for index, role := range user.Roles {
			log.Printf("\t\t\tRole %d: \n", index)
			log.Print("\t\t\t\tRole-Name: ", role.Name)
			log.Print("\t\t\t\tRole-permissions: ", role.Permissions)
		}
	}

}

func addRole(roleClient pb.RoleServiceClient, name pb.Name, permissions []pb.Permissions) (string, error) {
	role := &pb.Role{
		Name:        name,
		Permissions: permissions,
	}
	req := &pb.AddRoleRequest{
		Role: role,
	}

	res, err := roleClient.AddRole(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			return "", ErrAlreadyExists
		} else {
			return "", err
		}
	}

	return res.Name, nil
}

func fetchRole(roleClient pb.RoleServiceClient, roleName string) (*pb.Role, error) {
	req := &pb.FindRoleRequest{
		Name: roleName,
	}

	res, err := roleClient.FindRole(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return res.Role, nil
}

func fetchAllRoles(userClient pb.RoleServiceClient) ([]*pb.Role, error) {
	req := &pb.FindRolesRequest{}

	res, err := userClient.FindRoles(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return res.Roles, nil
}

func addUser(userClient pb.UserServiceClient, name, userAddress string, roles []*pb.Role) (string, error) {
	// user := sample.NewUser()
	user := &pb.User{
		Id:      uuid.New().String(),
		Name:    name,
		Address: userAddress,
		Roles:   roles,
	}

	// user.Id = "invalid-uuid"
	req := &pb.AddUserRequest{
		User: user,
	}

	res, err := userClient.AddUser(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			return "", ErrAlreadyExists
		} else {
			return "", err
		}
	}

	return res.Id, nil
}

func fetchUser(userClient pb.UserServiceClient, userID string) (*pb.User, error) {
	req := &pb.FindUserRequest{
		Id: userID,
	}

	res, err := userClient.FindUser(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return res.User, nil
}

func fetchAllUsers(userClient pb.UserServiceClient) ([]*pb.User, error) {
	req := &pb.FindUsersRequest{}

	res, err := userClient.FindUsers(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return res.Users, nil
}

func randomStringFromSet(a ...string) string {
	n := len(a)
	if n == 0 {
		return ""
	}
	return a[rand.Intn(n)]
}
