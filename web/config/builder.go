package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/labbsr0x/whisper-client/client"

	"github.com/labbsr0x/whisper/db"

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
	logLevel       = "log-level"
	scopesFilePath = "scopes-file-path"
	databaseURL    = "database-url"
	secretKey      = "secret-key"
	mailUser       = "mail-user"
	mailPassword   = "mail-password"
	mailIdentity   = "mail-identity"
	mailHost       = "mail-host"
	mailPort       = "mail-port"
)

// Flags define the fields that will be passed via cmd
type Flags struct {
	Port           string
	BaseUIPath     string
	LogLevel       string
	ScopesFilePath string
	HydraAdminURL  string
	HydraPublicURL string
	DatabaseURL    string
	SecretKey      string
	MailUser       string
	MailPassword   string
	MailIdentity   string
	MailHost       string
	MailPort       string
}

// WebBuilder defines the parametric information of a whisper server instance
type WebBuilder struct {
	*Flags
	Self               *client.WhisperClient
	GrantScopes        misc.GrantScopes
	UserCredentialsDAO db.UserCredentialsDAO
}

// AddFlags adds flags for Builder.
func AddFlags(flags *pflag.FlagSet) {
	flags.StringP(baseUIPath, "u", "", "Base path where the 'static' folder will be found with all the UI files")
	flags.StringP(port, "p", "7070", "[optional] Custom port for accessing Whisper's services. Defaults to 7070")
	flags.StringP(hydraAdminURL, "a", "", "Hydra Admin URL")
	flags.StringP(hydraPublicURL, "o", "", "Hydra Public URL")
	flags.StringP(logLevel, "l", "info", "[optional] Sets the Log Level to one of seven (trace, debug, info, warn, error, fatal, panic). Defaults to info")
	flags.StringP(scopesFilePath, "s", "", "Sets the path to the json file where the available scopes will be found")
	flags.StringP(databaseURL, "d", "", "Sets the database url where user credential data will be stored")
	flags.StringP(secretKey, "k", "", "Sets the secret key used to hash the stored passwords")
	flags.StringP(mailUser, "", "", "Sets the mail worker user")
	flags.StringP(mailPassword, "", "", "Sets the mail worker user's password")
	flags.StringP(mailIdentity, "", "", "Sets the mail worker user's identity")
	flags.StringP(mailHost, "", "", "Sets the mail worker host")
	flags.StringP(mailPort, "", "", "Sets the mail worker port")
}

// InitFromViper initializes the web server builder with properties retrieved from Viper.
func (b *WebBuilder) InitFromViper(v *viper.Viper) *WebBuilder {
	flags := new(Flags)
	flags.Port = v.GetString(port)
	flags.BaseUIPath = v.GetString(baseUIPath)
	flags.LogLevel = v.GetString(logLevel)
	flags.ScopesFilePath = v.GetString(scopesFilePath)
	flags.HydraAdminURL = v.GetString(hydraAdminURL)
	flags.HydraPublicURL = v.GetString(hydraPublicURL)
	flags.DatabaseURL = v.GetString(databaseURL)
	flags.SecretKey = v.GetString(secretKey)
	flags.MailUser = v.GetString(mailUser)
	flags.MailPassword = v.GetString(mailPassword)
	flags.MailIdentity = v.GetString(mailIdentity)
	flags.MailHost = v.GetString(mailHost)
	flags.MailPort = v.GetString(mailPort)

	flags.check()

	b.Flags = flags
	b.GrantScopes = b.getGrantScopesFromFile(flags.ScopesFilePath)
	b.Self = new(client.WhisperClient).InitFromParams(flags.HydraAdminURL, flags.HydraPublicURL, "whisper", "", b.GrantScopes.GetScopeListFromGrantScopeMap(), []string{})
	b.UserCredentialsDAO = new(db.DefaultUserCredentialsDAO).Init(b.DatabaseURL, b.SecretKey)

	logrus.Infof("GrantScopes: '%v'", b.GrantScopes)
	return b
}

func (flags *Flags) check() {
	logrus.Infof("Flags: '%v'", flags)

	haveEmptyRequiredFlags := flags.BaseUIPath == "" ||
		flags.HydraAdminURL == "" ||
		flags.HydraPublicURL == "" ||
		flags.ScopesFilePath == "" ||
		flags.SecretKey == "" ||
		flags.DatabaseURL == "" ||
		flags.MailUser == "" ||
		flags.MailPassword == "" ||
		flags.MailHost == "" ||
		flags.MailPort == ""

	if haveEmptyRequiredFlags {
		msg := fmt.Sprintf("The following flags cannot be empty:\n\n\t%v\n\t%v\n\t%v\n\t%v\n\t%v\n\t%v\n\t%v\n\t%v\n\t%v\n\t%v\n",
			baseUIPath,
			hydraAdminURL,
			hydraPublicURL,
			scopesFilePath,
			secretKey,
			databaseURL,
			mailUser,
			mailPassword,
			mailHost,
			mailPort)

		panic(msg)
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
