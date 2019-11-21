package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper/hydra"
	"github.com/labbsr0x/whisper/mail"

	"github.com/labbsr0x/whisper-client/client"
	"github.com/labbsr0x/whisper-client/config"

	"github.com/labbsr0x/whisper/misc"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	baseUIPath     = "base-ui-path"
	port           = "port"
	hydraAdminURL  = "hydra-admin-url"
	hydraPublicURL = "hydra-public-url"
	publicURL      = "public-url"
	logLevel       = "log-level"
	scopesFilePath = "scopes-file-path"
	databaseURL    = "database-url"
	secretKey      = "secret-key"
	mailUser       = "mail-user"
	mailPassword   = "mail-password"
	mailHost       = "mail-host"
	mailPort       = "mail-port"
	shutdownTime   = "shutdown-time"
)

// Flags define the fields that will be passed via cmd
type Flags struct {
	Port           string
	BaseUIPath     string
	LogLevel       string
	ScopesFilePath string
	HydraAdminURL  string
	HydraPublicURL string
	PublicURL      string
	DatabaseURL    string
	SecretKey      string
	MailUser       string
	MailPassword   string
	MailHost       string
	MailPort       string
	ShutdownTime   time.Duration
}

// WebBuilder defines the parametric information of a whisper server instance
type WebBuilder struct {
	*Flags
	Self        *client.WhisperClient
	HydraHelper hydra.Api
	GrantScopes misc.GrantScopes
	Outbox      chan<- mail.Mail
	DB          *gorm.DB
}

// AddFlags adds flags for Builder.
func AddFlags(flags *pflag.FlagSet) {
	flags.StringP(baseUIPath, "u", "", "Base path where the 'static' folder will be found with all the UI files")
	flags.StringP(port, "p", "7070", "[optional] Custom port for accessing Whisper's services. Defaults to 7070")
	flags.StringP(hydraAdminURL, "a", "", "Hydra Admin URL")
	flags.StringP(hydraPublicURL, "o", "", "Hydra Public URL")
	flags.StringP(publicURL, "", "", "Public URL for referencing in links")
	flags.StringP(logLevel, "l", "info", "[optional] Sets the Log Level to one of seven (trace, debug, info, warn, error, fatal, panic). Defaults to info")
	flags.StringP(scopesFilePath, "s", "", "Sets the path to the json file where the available scopes will be found")
	flags.StringP(databaseURL, "d", "", "Sets the database url where user credential data will be stored")
	flags.StringP(secretKey, "k", "", "Sets the secret key used to hash the stored passwords")
	flags.StringP(mailUser, "", "", "Sets the mail worker user")
	flags.StringP(mailPassword, "", "", "Sets the mail worker user's password")
	flags.StringP(mailHost, "", "", "Sets the mail worker host")
	flags.StringP(mailPort, "", "", "Sets the mail worker port")
	flags.StringP(shutdownTime, "t", "5", "[optional] Sets the Graceful Shutdown wait time (seconds). Defaults to 5")
}

// Init initializes the web server builder with properties retrieved from Viper.
func (b *WebBuilder) Init(v *viper.Viper, outbox chan<- mail.Mail) *WebBuilder {
	flags := new(Flags)
	flags.Port = v.GetString(port)
	flags.BaseUIPath = v.GetString(baseUIPath)
	flags.LogLevel = v.GetString(logLevel)
	flags.ScopesFilePath = v.GetString(scopesFilePath)
	flags.HydraAdminURL = v.GetString(hydraAdminURL)
	flags.HydraPublicURL = v.GetString(hydraPublicURL)
	flags.PublicURL = v.GetString(publicURL)
	flags.DatabaseURL = v.GetString(databaseURL)
	flags.SecretKey = v.GetString(secretKey)
	flags.MailUser = v.GetString(mailUser)
	flags.MailPassword = v.GetString(mailPassword)
	flags.MailHost = v.GetString(mailHost)
	flags.MailPort = v.GetString(mailPort)
	flags.ShutdownTime = v.GetDuration(shutdownTime)

	flags.check()

	b.Flags = flags
	b.Outbox = outbox
	b.GrantScopes = b.getGrantScopesFromFile(flags.ScopesFilePath)
	b.HydraHelper = new(hydra.DefaultHydraHelper).Init(b.HydraAdminURL)
	b.DB = b.initDB()

	hydraAdminURI, err := url.Parse(flags.HydraAdminURL)
	gohtypes.PanicIfError("Invalid hydra admin url", 500, err)
	hydraPublicURI, err := url.Parse(flags.HydraPublicURL)
	gohtypes.PanicIfError("Invalid hydra public url", 500, err)

	b.Self = new(client.WhisperClient).InitFromConfig(&config.Config{
		ClientID:       "whisper",
		ClientSecret:   "",
		WhisperURL:     nil,
		HydraAdminURL:  hydraAdminURI,
		HydraPublicURL: hydraPublicURI,
		Scopes:         b.GrantScopes.GetScopeListFromGrantScopeMap(),
	})

	logrus.Infof("GrantScopes: '%v'", b.GrantScopes)
	return b
}

func (flags *Flags) check() {
	logrus.Infof("Flags: '%v'", flags)

	requiredFlags := []struct{value string; name string}{
		{flags.BaseUIPath, baseUIPath},
		{flags.HydraAdminURL, hydraAdminURL},
		{flags.HydraPublicURL, hydraPublicURL},
		{flags.PublicURL, publicURL},
		{flags.ScopesFilePath, scopesFilePath},
		{flags.SecretKey, secretKey},
		{flags.DatabaseURL, databaseURL},
		{flags.MailUser, mailUser},
		{flags.MailPassword, mailPassword},
		{flags.MailHost, mailHost},
		{flags.MailPort, mailPort},
	}

	var errMsg string

	for _, flag := range requiredFlags {
		if flag.value == "" {
			errMsg += fmt.Sprintf("\n\t%v", flag.name)
		}
	}

	if errMsg != "" {
		errMsg = "The following flags are missing: " + errMsg
		panic(errMsg)
	}
}

// getGrantScopesFromFile reads into memory the json scopes file
func (b *WebBuilder) getGrantScopesFromFile(scopesFilePath string) misc.GrantScopes {
	jsonFile, err := os.Open(scopesFilePath)
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	var grantScopes misc.GrantScopes

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(bytes, &grantScopes)
	if err != nil {
		panic(err.Error())
	}

	return grantScopes
}

// initDB opens a connection with the database
func (b *WebBuilder) initDB() *gorm.DB {
	dbURL := strings.Replace(b.DatabaseURL, "mysql://", "", 1)
	dbc, err := gorm.Open("mysql", dbURL)
	gohtypes.PanicIfError("Unable to open db", http.StatusInternalServerError, err)

	dbc.DB().SetMaxIdleConns(4)
	dbc.DB().SetMaxOpenConns(20)

	return dbc.LogMode(true)

}
