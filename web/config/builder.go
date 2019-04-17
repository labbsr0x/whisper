package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/labbsr0x/whisper/misc"

	"github.com/labbsr0x/whisper-client/hydra"
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
	clientID       = "client-id"
	clientSecret   = "client-secret"
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
	ClientID       string
	ClientSecret   string
}

// WebBuilder defines the parametric information of a whisper server instance
type WebBuilder struct {
	*Flags
	HydraClient *hydra.Client
	GrantScopes misc.GrantScopes
}

// AddFlags adds flags for Builder.
func AddFlags(flags *pflag.FlagSet) {
	flags.String(baseUIPath, "", "Base path where the 'static' folder will be found with all the UI files")
	flags.String(port, "7070", "[optional] Custom port for accessing Whisper's services. Defaults to 7070")
	flags.String(hydraAdminURL, "", "Hydra Admin URL")
	flags.String(hydraPublicURL, "", "Hydra Public URL")
	flags.String(logLevel, "info", "[optional] Sets the Log Level to one of seven (trace, debug, info, warn, error, fatal, panic). Defaults to info")
	flags.String(scopesFilePath, "", "Sets the path to the json file where the available scopes will be found")
	flags.String(databaseURL, "", "Sets the database url where user credential data will be stored")
	flags.String(clientID, "", "Sets the oauth2 client id")
	flags.String(clientSecret, "", "Sets the oauth2 client secret")
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
	flags.ClientID = v.GetString(clientID)
	flags.ClientSecret = v.GetString(clientSecret)
	flags.DatabaseURL = v.GetString(databaseURL)

	flags.check()

	b.Flags = flags
	b.GrantScopes = b.getGrantScopesFromFile(flags.ScopesFilePath)
	b.HydraClient = new(hydra.Client).Init(flags.HydraAdminURL, flags.HydraPublicURL, flags.ClientID, flags.ClientSecret, b.GrantScopes.GetScopeListFromGrantScopeMap(), []string{})

	logrus.Infof("Run config: %v", misc.GetJSONStr(b))
	return b
}

func (flags *Flags) check() {

	if flags.BaseUIPath == "" || flags.HydraAdminURL == "" || flags.HydraPublicURL == "" || flags.ScopesFilePath == "" || flags.ClientID == "" || flags.ClientSecret == "" {
		panic("base-ui-path, hydra-admin-url, hydra-public-url, scopes-file-path, client-id and client-secret cannot be empty")
	}

	if len(flags.ClientSecret) < 6 {
		panic("client-secret must be at least 6 characters long")
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
	bytes, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(bytes, &grantScopes)

	return grantScopes
}
