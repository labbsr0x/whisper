package db

import (
	"github.com/labbsr0x/whisper/mail"
	"net/http"
	"time"

	"github.com/labbsr0x/whisper/misc"

	"github.com/google/uuid"

	"github.com/labbsr0x/goh/gohtypes"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// UserCredential holds the information from a user credential
type UserCredential struct {
	ID             string `gorm:"primary_key;not null;"`
	Username       string `gorm:"unique_index;not null;"`
	Email          string `gorm:"unique_index;not null;"`
	Password       string `gorm:"not null;"`
	Salt           string `gorm:"not null;"`
	EmailValidated bool   `gorm:"not null;"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// BeforeCreate will set a UUID rather than numeric ID.
func (user *UserCredential) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

// UserCredentialsDAO defines the methods that can be performed
type UserCredentialsDAO interface {
	Init(secretKey, baseUIPath, publicAddressURL string, outbox chan<- mail.Mail, db *gorm.DB) UserCredentialsDAO
	CreateUserCredential(username, password, email string) (string, error)
	UpdateUserCredential(username, email, password string) error
	GetUserCredential(username string) (UserCredential, error)
	CheckCredentials(username, password string) UserCredential
	ValidateUserCredentialEmail(username string) error
}

// DefaultUserCredentialsDAO a default UserCredentialsDAO interface implementation
type DefaultUserCredentialsDAO struct {
	db               *gorm.DB
	outbox           chan<- mail.Mail
	secretKey        string
	baseUIPath       string
	publicAddressURL string
}

// InitFromWebBuilder initializes a default user credentials DAO from web builder
func (dao *DefaultUserCredentialsDAO) Init(secretKey, baseUIPath, publicAddressURL string, outbox chan<- mail.Mail, db *gorm.DB) UserCredentialsDAO {
	dao.secretKey = secretKey
	dao.outbox = outbox
	dao.db = db
	dao.baseUIPath = baseUIPath
	dao.publicAddressURL = publicAddressURL

	err := dao.db.AutoMigrate(&UserCredential{}).Error
	gohtypes.PanicIfError("Not possible to migrate db", http.StatusInternalServerError, err)

	return dao
}

// CreateUserCredential creates a user
func (dao *DefaultUserCredentialsDAO) CreateUserCredential(username, password, email string) (string, error) {
	var users []UserCredential

	if res := dao.db.Model(&UserCredential{}).Where("username = ?", username).Or("email = ?", email).Find(&users); res.Error != nil {
		return "", res.Error
	}

	for _, user := range users {
		if user.Username == username {
			gohtypes.Panic("Username already taken", http.StatusConflict)
		}

		if user.Email == email {
			gohtypes.Panic("Email already taken", http.StatusConflict)
		}
	}

	salt := misc.GenerateSalt()
	hPassword := misc.GetEncryptedPassword(dao.secretKey, password, salt)
	userCredential := UserCredential{
		Username:       username,
		Password:       hPassword,
		Email:          email,
		Salt:           salt,
		EmailValidated: false,
	}

	if res := dao.db.Create(&userCredential); res.Error != nil {
		return "", res.Error
	}

	return userCredential.ID, nil
}

func (dao *DefaultUserCredentialsDAO) ValidateUserCredentialEmail(username string) error {
	userCredential, err := dao.GetUserCredential(username)
	if err != nil {
		return err
	}

	userCredential.EmailValidated = true

	return dao.db.Save(userCredential).Error
}

// UpdateUserCredential updates a user
func (dao *DefaultUserCredentialsDAO) UpdateUserCredential(username, email, password string) error {
	userCredential := UserCredential{}

	err := dao.db.Where("username = ?", username).First(&userCredential).Error
	gohtypes.PanicIfError("Unable to retrieve user", http.StatusInternalServerError, err)

	if hNewPassword := misc.GetEncryptedPassword(dao.secretKey, password, userCredential.Salt); hNewPassword != userCredential.Password {
		salt := misc.GenerateSalt()
		hPassword := misc.GetEncryptedPassword(dao.secretKey, password, salt)

		userCredential.Password = hPassword
		userCredential.Salt = salt
	}

	if email != userCredential.Email {
		userCredential.Email = email
		userCredential.EmailValidated = false

		dao.outbox <- mail.GetEmailConfirmationMail(dao.baseUIPath, dao.secretKey, dao.publicAddressURL, username, email, "")
	}

	return dao.db.Save(userCredential).Error
}

// GetUserCredential gets an user credential
func (dao *DefaultUserCredentialsDAO) GetUserCredential(username string) (userCredential UserCredential, err error) {
	err = dao.db.Where("username = ?", username).First(&userCredential).Error
	return
}

// CheckCredentials verifies if the informed credentials are valid
func (dao *DefaultUserCredentialsDAO) CheckCredentials(username, password string) UserCredential {
	userCredential, err := dao.GetUserCredential(username)

	if err != nil {
		gohtypes.PanicIfError("Unable to authenticate user", http.StatusInternalServerError, err)
	}

	hPassword := misc.GetEncryptedPassword(dao.secretKey, password, userCredential.Salt)

	if hPassword != userCredential.Password {
		gohtypes.Panic("Incorrect password", http.StatusUnauthorized)
	}

	return userCredential
}
