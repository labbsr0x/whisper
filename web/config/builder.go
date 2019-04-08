package config

import (
	"github.com/abilioesteves/whisper/misc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	baseUIPath         = "base-ui-path"
	port               = "port"
	hydraAdminEndpoint = "hydra-admin-endpoint"
	logLevel           = "log-level"
	scopesFilePath     = "scopes-file-path"
)

// GrantScope defines the structure of a grant scope
type GrantScope struct {
	Description string
	Details     string
	Scope       string
}

// Flags define the fields that will be passed via cmd
type Flags struct {
	Port               string
	BaseUIPath         string
	LogLevel           string
	ScopesFilePath     string
	HydraAdminEndpoint string
}

// WebBuilder defines the parametric information of a whisper server instance
type WebBuilder struct {
	*Flags
	HydraClient *misc.HydraClient
	GrantScopes map[string]GrantScope
}

// AddFlags adds flags for Builder.
func AddFlags(flags *pflag.FlagSet) {
	flags.String(baseUIPath, "", "Base path where the 'static' folder will be found with all the UI files")
	flags.String(port, "7070", "Custom port for accessing Whisper's services")
	flags.String(hydraAdminEndpoint, "", "Hydra Admin Enpoint")
	flags.String(logLevel, "info", "Sets the Log Level to one of seven (trace, debug, info, warn, error, fatal, panic)")
	flags.String(scopesFilePath, "", "Sets the path to the json file where the available scopes will be found")
}

// InitFromViper initializes the web server builder with properties retrieved from Viper.
func (b *WebBuilder) InitFromViper(v *viper.Viper) *WebBuilder {
	flags := new(Flags)
	flags.Port = v.GetString(port)
	flags.BaseUIPath = v.GetString(baseUIPath)
	flags.LogLevel = v.GetString(logLevel)
	flags.ScopesFilePath = v.GetString(scopesFilePath)
	flags.HydraAdminEndpoint = v.GetString(hydraAdminEndpoint)

	flags.check()

	b.Flags = flags
	b.HydraClient = new(misc.HydraClient).Init(flags.HydraAdminEndpoint)
	b.GrantScopes = b.getGrantScopesFromFile(flags.ScopesFilePath)

	logrus.Infof("Run config: %v", b)
	return b
}

func (flags *Flags) check() {
	if flags.BaseUIPath == "" || flags.HydraAdminEndpoint == "" || flags.ScopesFilePath == "" {
		panic("base-ui-path, hydra-admin-endpoint and scopes-file-path cannot be empty")
	}
}

// getGrantScopesFromFile reads into memory the json scopes file
func (b *WebBuilder) getGrantScopesFromFile(scopesFilePath string) map[string]GrantScope {
	return map[string]GrantScope{
		"openid": GrantScope{
			Description: "Access to your personal data",
			Scope:       "openid",
			Details:     "Provides access to personal data such as: email, name etc",
		},
		"offline": GrantScope{
			Description: "Always Sign in",
			Scope:       "offline",
			Details:     "Provides the possibility for the app to be always signed in to your account",
		},
	}
}
