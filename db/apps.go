package db

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper/misc"
	"github.com/labbsr0x/whisper/web/api/types"
)

// AppData holds the information from an App
type AppData struct {
	ID           string `gorm:"primary_key;not null;"`
	ClientID     string `gorm:"unique_index;not null;"`
	ClientSecret string `gorm:"not null;"`
	ClientName   string `gorm:"not null;"`
	Type         string `gorm:"not null;"`
	Owner        string `gorm:"not null;"`

	Address           string
	LoginRedirectURL  string
	LogoutRedirectURL string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// AppDAO defines the interface of the AppDAO
type AppDAO interface {
	Init(db *gorm.DB) AppDAO
	InsertAppData(payload types.AddAppInitialRequestPayload, owner string) (types.App, error)
	UpdateAppData(payload types.UpdateAppRequestPayload, owner string) (types.App, error)
	ListApps(owner string) ([]types.App, error)
}

// DefaultAppDAO defines the default implementation of the AppDAO interface
type DefaultAppDAO struct {
	db *gorm.DB
}

// BeforeCreate will set a UUID rather than numeric ID.
func (app *AppData) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

// Init initializes the dao
func (dao *DefaultAppDAO) Init(db *gorm.DB) AppDAO {
	dao.db = db

	err := dao.db.AutoMigrate(&AppData{}).Error
	gohtypes.PanicIfError("Not possible to migrate db", http.StatusInternalServerError, err)

	return dao
}

// InsertAppData inserts an app to the db
func (dao *DefaultAppDAO) InsertAppData(payload types.AddAppInitialRequestPayload, owner string) (types.App, error) {
	app := AppData{
		ClientID:          uuid.New().String(),
		ClientSecret:      misc.GetEncryptedPassword(uuid.New().String(), "", ""),
		ClientName:        payload.Name,
		Address:           "",
		LoginRedirectURL:  payload.LoginRedirectURL,
		LogoutRedirectURL: payload.LogoutRedirectURL,
		Owner:             owner,
		Type:              payload.Type,
	}

	if res := dao.db.Create(&app); res.Error != nil {
		return types.App{}, res.Error
	}

	return types.App{
		ClientID:          app.ClientID,
		ClientName:        app.ClientName,
		ClientSecret:      app.ClientSecret,
		Address:           app.Address,
		LoginRedirectURL:  app.LoginRedirectURL,
		LogoutRedirectURL: app.LogoutRedirectURL,
		Type:              app.Type,
	}, nil
}

// UpdateAppData updates the stored app data
func (dao *DefaultAppDAO) UpdateAppData(payload types.UpdateAppRequestPayload, owner string) (types.App, error) {
	var appData AppData
	if res := dao.db.Model(&types.App{}).Where("clientid = ? AND owner = ?", payload.ClientID, owner).First(&appData); res.Error != nil {
		return types.App{}, res.Error
	}

	appData.ClientSecret = payload.Secret
	appData.ClientName = payload.Name
	appData.Address = payload.Address
	appData.LoginRedirectURL = payload.LoginRedirectURL
	appData.LogoutRedirectURL = payload.LogoutRedirectURL

	err := dao.db.Save(&appData).Error
	if err != nil {
		return types.App{}, err
	}
	return types.App{
		ClientID:          appData.ClientID,
		ClientName:        appData.ClientName,
		ClientSecret:      appData.ClientSecret,
		Address:           appData.Address,
		LoginRedirectURL:  appData.LoginRedirectURL,
		LogoutRedirectURL: appData.LogoutRedirectURL,
		Type:              appData.Type,
	}, nil
}

// ListApps lists the apps stored in the db
func (dao *DefaultAppDAO) ListApps(owner string) ([]types.App, error) {
	var apps []types.App
	if res := dao.db.Model(&types.App{}).Where("owner = ?", owner).Find(&apps); res.Error != nil {
		return nil, res.Error
	}

	return apps, nil
}
