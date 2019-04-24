package db

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"

	"github.com/labbsr0x/goh/gohtypes"

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
	GetUserCredential(username string) (UserCredential, error)
	InitFromDatabaseURL(dbURL string) UserCredentialsDAO
}

// DefaultUserCredentialsDAO a default UserCredentialsDAO interface implementation
type DefaultUserCredentialsDAO struct {
	DatabaseURL string
}

// InitFromDatabaseURL initializes a defualt user credentials DAO from web builder
func (dao *DefaultUserCredentialsDAO) InitFromDatabaseURL(dbURL string) UserCredentialsDAO {
	u, err := url.Parse(dbURL)
	gohtypes.PanicIfError("Unable to parse db url", 500, err)
	dao.DatabaseURL = strings.Replace(u.String(), u.Scheme+"://", "", 1)

	gohtypes.PanicIfError("Not possible to migrate db", 500, dao.migrate())

	return dao
}

// migrate initializes a migration routine to synchronize db and model
func (dao *DefaultUserCredentialsDAO) migrate() error {
	db, err := gorm.Open("mysql", dao.DatabaseURL)
	if err == nil {
		defer db.Close()
		db.AutoMigrate(&UserCredential{})
	}
	logrus.Error(err)
	return err
}

// CreateUserCredential creates a user
func (dao *DefaultUserCredentialsDAO) CreateUserCredential(username, password, email string) (string, error) {
	db, err := gorm.Open("mysql", dao.DatabaseURL)
	if err == nil {
		defer db.Close()
		userCredential := UserCredential{ID: uuid.New().String(), Username: username, Password: password, Email: email, Salt: password}
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

// GetUserCredential gets an user credential
func (dao *DefaultUserCredentialsDAO) GetUserCredential(username string) (UserCredential, error) {
	userCredential := UserCredential{}
	db, err := gorm.Open("mysql", dao.DatabaseURL)
	if err == nil {
		defer db.Close()

		db.Where("username = ?", username).First(&userCredential)
	}
	return userCredential, err
}
