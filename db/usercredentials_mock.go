package db

import (
	"github.com/google/uuid"
)

// CreateUserCredential creates a user
func CreateUserCredential(username, password, clientID string) (string, error) {
	return "", nil
}

// UpdateUserCredential updates a user
func UpdateUserCredential(userID uuid.UUID, username, password string) error {
	return nil
}

// DeleteUserCredential deletes the user
func DeleteUserCredential(userID uuid.UUID) error {
	return nil
}
