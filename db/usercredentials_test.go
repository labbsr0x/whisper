package db

import (
	"time"

	"github.com/google/uuid"
)

// UserCredential structures the information that defines a user credential
type UserCredential struct {
	ID          string    `db:"id"`
	ClientID    string    `db:"client_id"`
	Username    string    `db:"username"`
	Password    string    `db:"password"`
	DateCreated time.Time `db:"date_created"`
}

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
