package db

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/labbsr0x/goh/gohtypes"
	apiTypes "github.com/labbsr0x/whisper/web/api/types"
)

// App holds the information from an App
type App struct {
	ID                string `gorm:"primary_key;not null;"`
	ClientID          string `gorm:"unique_index;not null;"`
	ClientSecret      string
	ClientName        string `gorm:"not null;"`
	ClientURL         string `gorm:"not null;"`
	LoginRedirectURL  string `gorm:"not null;"`
	LogoutRedirectURL string `gorm:"not null;"`
	GrantTypes        string `gorm:"not null;"`
	Owner             string `gorm:"not null;"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// AppDAO defines the interface of the AppDAO
type AppDAO interface {
	Init(db *gorm.DB) AppDAO
	CreateApp(payload apiTypes.AddAppRequestPayload, owner string) error
	ListApps(owner string) ([]App, error)
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

	err := dao.db.AutoMigrate(&App{}).Error
	gohtypes.PanicIfError("Not possible to migrate db", http.StatusInternalServerError, err)

	return dao
}

// CreateApp inserts an app to the db
func (dao *DefaultAppDAO) CreateApp(payload apiTypes.AddAppRequestPayload, owner string) error {
	var apps []App
	if res := dao.db.Model(&App{}).Where("clientId = ?", payload.ID).Find(&apps); res.Error != nil {
		return res.Error
	}

	app := App{
		ClientID:          payload.ID,
		ClientSecret:      payload.Secret,
		ClientName:        payload.Name,
		ClientURL:         payload.URL,
		LoginRedirectURL:  payload.LoginRedirectURL,
		LogoutRedirectURL: payload.LogoutRedirectURL,
		GrantTypes:        strings.Join(payload.GrantTypes, ","),
		Owner:             owner,
	}

	if res := dao.db.Create(&app); res.Error != nil {
		return res.Error
	}

	return nil
}

// ListApps lists the apps stored in the db
func (dao *DefaultAppDAO) ListApps(owner string) ([]App, error) {
	var apps []App
	if res := dao.db.Model(&App{}).Where("owner = ?", owner).Find(&apps); res.Error != nil {
		return nil, res.Error
	}

	return apps, nil
}
