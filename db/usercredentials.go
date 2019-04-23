package db

import (
	"fmt"
	"time"

	"github.com/labbsr0x/whisper/web/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// UserCredential holds the information from a user credential
type UserCredential struct {
	ID        string `gorm:"primary_key;not null;"`
	Username  string `gorm:"unique_index;not null;"`
	Email     string `gorm:"index"`
	Password  string `gorm:"not null;"`
	Salt      string `gorm:"not null;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserCredentialsDAO defines the methods that can be performed
type UserCredentialsDAO interface {
	CreateUserCredential(username, password, email string) (string, error)
	UpdateUserCredential(userID string, email, password string) error
	DeleteUserCredential(userID string) error
}

// DefaultUserCredentialsDAO a default UserCredentialsDAO interface implementation
type DefaultUserCredentialsDAO struct {
	DatabaseURL string
}

// InitFromWebBuilder initializes a defualt user credentials DAO from web builder
func (dao *DefaultUserCredentialsDAO) InitFromWebBuilder(builder *config.WebBuilder) *DefaultUserCredentialsDAO {
	dao.DatabaseURL = builder.DatabaseURL
	return dao
}

// CreateUserCredential creates a user
func (dao *DefaultUserCredentialsDAO) CreateUserCredential(username, password, email string) (string, error) {
	db, err := gorm.Open("mysql", dao.DatabaseURL)
	if err == nil {
		userCredential := UserCredential{Username: username, Password: password, Email: email, Salt: password}
		db.NewRecord(userCredential)

		db.Create(&userCredential)

		if !db.NewRecord(userCredential) {
			return userCredential.ID, nil
		}

		err = fmt.Errorf("Unable to create an user credential: %v", db.GetErrors())
	}
	return "", err
}

// UpdateUserCredential updates a user
func (dao *DefaultUserCredentialsDAO) UpdateUserCredential(userID string, email, password string) error {
	return nil
}

// DeleteUserCredential deletes the user
func (dao *DefaultUserCredentialsDAO) DeleteUserCredential(userID string) error {
	return nil
}
