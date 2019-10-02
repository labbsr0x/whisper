package db

import (
	"github.com/labbsr0x/whisper/mail"
	"github.com/labbsr0x/whisper/resources"
	"net/http"
	"strings"
	"time"

	"github.com/labbsr0x/whisper/misc"

	"github.com/google/uuid"

	"github.com/labbsr0x/goh/gohtypes"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// UserCredential holds the information from a user credential
type UserCredential struct {
	ID            string `gorm:"primary_key;not null;"`
	Username      string `gorm:"unique_index;not null;"`
	Email         string `gorm:"unique_index;not null;"`
	Password      string `gorm:"not null;"`
	Salt          string `gorm:"not null;"`
	Authenticated bool   `gorm:"not null;"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// BeforeCreate will set a UUID rather than numeric ID.
func (user *UserCredential) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

// UserCredentialsDAO defines the methods that can be performed
type UserCredentialsDAO interface {
	Init(dbURL, secretKey string) UserCredentialsDAO
	CreateUserCredential(username, password, email string) (string, error)
	UpdateUserCredential(username, email, password string) error
	GetUserCredential(username string) (UserCredential, error)
	CheckCredentials(username, password, challenge string)
	AuthenticateUserCredential(username string)
}

// DefaultUserCredentialsDAO a default UserCredentialsDAO interface implementation
type DefaultUserCredentialsDAO struct {
	DatabaseURL string
	SecretKey   string
}

// Init initializes a default user credentials DAO from web builder
func (dao *DefaultUserCredentialsDAO) Init(dbURL, secretKey string) UserCredentialsDAO {
	dao.SecretKey = secretKey
	dao.DatabaseURL = strings.Replace(dbURL, "mysql://", "", 1)

	gohtypes.PanicIfError("Not possible to migrate db", http.StatusInternalServerError, dao.migrate())

	return dao
}

// migrate initializes a migration routine to synchronize db and model
func (dao *DefaultUserCredentialsDAO) migrate() error {
	db, err := gorm.Open("mysql", dao.DatabaseURL)
	if err == nil {
		defer db.Close()
		db.LogMode(true)
		err = db.AutoMigrate(&UserCredential{}).Error
	}
	return err
}

// CreateUserCredential creates a user
func (dao *DefaultUserCredentialsDAO) CreateUserCredential(username, password, email string) (string, error) {
	db, err := gorm.Open("mysql", dao.DatabaseURL)
	if err != nil {
		return "", err
	}

	defer db.Close()

	var users []UserCredential

	if res := db.Model(&UserCredential{}).Where("username = ?", username).Or("email = ?", email).Find(&users); res.Error != nil {
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
	hPassword := misc.GetEncryptedPassword(dao.SecretKey, password, salt)
	userCredential := UserCredential{
		Username:      username,
		Password:      hPassword,
		Email:         email,
		Salt:          salt,
		Authenticated: false,
	}

	if res := db.Create(&userCredential); res.Error != nil {
		return "", res.Error
	}

	return userCredential.ID, nil
}

func (dao *DefaultUserCredentialsDAO) AuthenticateUserCredential(username string) {
	userCredential, err := dao.GetUserCredential(username)
	gohtypes.PanicIfError("Unable to retrieve user", http.StatusInternalServerError, err)

	userCredential.Authenticated = true

	db, err := gorm.Open("mysql", dao.DatabaseURL)
	gohtypes.PanicIfError("Unable to connect with database", http.StatusInternalServerError, err)

	if db := db.Save(userCredential); db.Error != nil {
		gohtypes.Panic("Unable to authenticate user", http.StatusInternalServerError)
	}
}

// UpdateUserCredential updates a user
func (dao *DefaultUserCredentialsDAO) UpdateUserCredential(username, email, password string) error {
	db, err := gorm.Open("mysql", dao.DatabaseURL)
	if err == nil {
		defer db.Close()

		salt := misc.GenerateSalt()
		hPassword := misc.GetEncryptedPassword(dao.SecretKey, password, salt)

		userCredential := UserCredential{}
		db.Where("username = ?", username).First(&userCredential)

		userCredential.Password = hPassword
		userCredential.Salt = salt
		userCredential.Email = email

		db = db.Save(userCredential)
		err = db.Error
	}
	return err
}

// GetUserCredential gets an user credential
func (dao *DefaultUserCredentialsDAO) GetUserCredential(username string) (UserCredential, error) {
	userCredential := UserCredential{}
	db, err := gorm.Open("mysql", dao.DatabaseURL)
	if err == nil {
		defer db.Close()
		db = db.Where("username = ?", username).First(&userCredential)
		err = db.Error
	}
	return userCredential, err
}

// CheckCredentials verifies if the informed credentials are valid
func (dao *DefaultUserCredentialsDAO) CheckCredentials(username, password, challenge string) {
	userCredential, err := dao.GetUserCredential(username)

	if err != nil {
		gohtypes.PanicIfError("Unable to authenticate user", http.StatusInternalServerError, err)
	}

	hPassword := misc.GetEncryptedPassword(dao.SecretKey, password, userCredential.Salt)

	if hPassword != userCredential.Password {
		gohtypes.Panic("Incorrect password", http.StatusUnauthorized)
	}

	if !userCredential.Authenticated {
		to, content := misc.GetEmailConfirmationMail(userCredential.Username, userCredential.Email, challenge)
		resources.Outbox <- mail.Mail{To: to, Content: content}
		gohtypes.Panic("This account email is not authenticated, an email was sent to you confirm your email", http.StatusUnauthorized)
	}
}
