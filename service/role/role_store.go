package service

import (
	"errors"
	"sync"

	"ashish/user-mgmt/pb"
)

var ErrAlreadyExists = errors.New("resource already exists")
var ErrNotFound = errors.New("resource not found")

type RoleStore interface {
	Save(role *pb.Role) error
	Find(id string) (*pb.Role, error)
	FindAll() (map[string]*pb.Role, error)
}

type InMemoryRoleStore struct {
	mutex sync.RWMutex // Since there will be multiple concurrent requests to save role, hence we need a READ-WRITE mutex.
	data  map[string]*pb.Role
}

func NewInMemoryRoleStore() *InMemoryRoleStore {
	return &InMemoryRoleStore{
		data: make(map[string]*pb.Role),
	}
}

func (roleStore *InMemoryRoleStore) Save(role *pb.Role) error {
	roleStore.mutex.Lock()
	defer roleStore.mutex.Unlock()

	if roleStore.data[role.Name.String()] != nil {
		return ErrAlreadyExists
	}

	roleStore.data[role.Name.String()] = role
	return nil
}

func (roleStore *InMemoryRoleStore) Find(id string) (*pb.Role, error) {
	roleStore.mutex.RLock()
	defer roleStore.mutex.RUnlock()

	role := roleStore.data[id]
	if role == nil {
		return nil, ErrNotFound
	}

	return role, nil
}

func (roleStore *InMemoryRoleStore) FindAll() (map[string]*pb.Role, error) {
	roleStore.mutex.RLock()
	defer roleStore.mutex.RUnlock()

	return roleStore.data, nil
}
