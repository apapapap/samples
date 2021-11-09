package service

import (
	"errors"
	"sync"

	"ashish/user-mgmt/pb"
)

var ErrAlreadyExists = errors.New("user already exists")
var ErrNotFound = errors.New("user not found")

type UserStore interface {
	Save(user *pb.User) error
	Find(id string) (*pb.User, error)
	FindAll() (map[string]*pb.User, error)
}

type InMemoryUserStore struct {
	mutex sync.RWMutex // Since there will be multiple concurrent requests to save user, hence we need a READ-WRITE mutex.
	data  map[string]*pb.User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		data: make(map[string]*pb.User),
	}
}

func (store *InMemoryUserStore) Save(user *pb.User) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[user.Id] != nil {
		return ErrAlreadyExists
	}

	store.data[user.Id] = user
	return nil
}

func (store *InMemoryUserStore) Find(id string) (*pb.User, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	user := store.data[id]
	if user == nil {
		return nil, ErrNotFound
	}

	return user, nil
}

func (store *InMemoryUserStore) FindAll() (map[string]*pb.User, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	return store.data, nil
}
