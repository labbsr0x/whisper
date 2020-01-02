package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// App holds the information from an App
type App struct {
	ID                string `gorm:"primary_key;not null;"`
	Secret            string
	Name              string `gorm:"unique_index;not null;"`
	URL               string `gorm:"not null;"`
	LoginRedirectURL  string `gorm:"not null;"`
	LogoutRedirectURL string `gorm:"not null;"`
	GrantTypes        []string
	Scopes            []string
	Owner             string `gorm:"not null;"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// AppDAO defines the interface of the AppDAO
type AppDAO interface {
	Init(db *gorm.DB) AppDAO
}

// DefaultAppDAO defines the default implementation of the AppDAO interface
type DefaultAppDAO struct {
	db *gorm.DB
}

// BeforeCreate will set a UUID rather than numeric ID.
func (app *App) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

// Init initializes the dao
func (dao *DefaultAppDAO) Init(db *gorm.DB) AppDAO {
	dao.db = db

	return dao
}
