package bl

import (
	"time"

	"github.com/google/uuid"
)

// User structures the information that defines a user
type User struct {
	ID          uuid.UUID `db:"id"`
	ClientID    string    `db:"client_id"`
	Username    string    `db:"username"`
	Password    string    `db:"password"`
	DateCreated time.Time `db:"date_created"`
}

// CreateUser creates a user
func CreateUser(username, password string, clientID uuid.UUID) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

// UpdateUser updates a user
func UpdateUser(userID uuid.UUID, username, password string) error {
	return nil
}

// DeleteUser deletes the user
func DeleteUser(userID uuid.UUID) error {
	return nil
}
