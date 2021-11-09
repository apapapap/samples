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

var ErrAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

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

	log.Print("1. User Add")
	name := randomStringFromSet("Ashish", "Pratik", "Rajat")
	userAddress := randomStringFromSet("Pune", "Mumbai", "Nashik")
	userID, err := addUser(userClient, name, userAddress)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			log.Print("user already exists")
		} else {
			log.Print("cannot create user: ", err)
		}
	} else {
		log.Print("user created with id: ", userID)
	}

	log.Print("2. List all users")
	users, nil := fetchAll(userClient)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			log.Print("user not found")
		} else {
			log.Print("cannot find user: ", err)
		}
	} else {
		for index, user := range users {
			log.Printf("	User: %d\n", index+1)
			log.Printf("		ID: %s\n", user.Id)
			log.Print("		Name: ", user.Name)
			log.Print("		Address: ", user.Address)
		}
	}

	log.Print("3. User Fetch")
	user, err := fetchUser(userClient, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			log.Print("user not found")
		} else {
			log.Print("cannot find user: ", err)
		}
	} else {
		log.Print("user found")
		log.Print("		ID: ", user.Id)
		log.Print("		Name: ", user.Name)
		log.Print("		Address: ", user.Address)
	}

}

func addUser(userClient pb.UserServiceClient, name, userAddress string) (string, error) {
	// user := sample.NewUser()
	user := &pb.User{
		Id:      uuid.New().String(),
		Name:    name,
		Address: userAddress,
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
			return res.User, ErrUserNotFound
		} else {
			return res.User, err
		}
	}
	return res.User, nil
}

func fetchAll(userClient pb.UserServiceClient) ([]*pb.User, error) {
	req := &pb.FindUsersRequest{}

	res, err := userClient.FindUsers(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return res.Users, ErrUserNotFound
		} else {
			return res.Users, err
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
