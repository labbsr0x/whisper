package db

import (
	"github.com/google/uuid"
)

// UserCredential holds the information from a user credential
type UserCredential struct {
}

// UserCredentialsDAO defines the methods that can be performed
type UserCredentialsDAO interface {
	CreateUserCredential(username, password, clientID string) (string, error)
	UpdateUserCredential(userID uuid.UUID, username, password string) error
	DeleteUserCredential(userID uuid.UUID) error
}

// DefaultUserCredentialsDAO a default UserCredentialsDAO interface implementation
type DefaultUserCredentialsDAO struct {
}

// CreateUserCredential creates a user
func (dao *DefaultUserCredentialsDAO) CreateUserCredential(username, password, clientID string) (string, error) {
	return "", nil
}

// UpdateUserCredential updates a user
func (dao *DefaultUserCredentialsDAO) UpdateUserCredential(userID uuid.UUID, username, password string) error {
	return nil
}

// DeleteUserCredential deletes the user
func (dao *DefaultUserCredentialsDAO) DeleteUserCredential(userID uuid.UUID) error {
	return nil
}
