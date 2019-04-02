package web

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	baseUIPath         = "base-ui-path"
	port               = "port"
	hydraAdminEndpoint = "hydra-admin-endpoint"
	logLevel           = "log-level"
)

// AddFlags adds flags for Builder.
func AddFlags(flags *pflag.FlagSet) {
	flags.String(baseUIPath, "", "Base path where the 'static' folder will be found with all the UI files")
	flags.String(port, "7070", "Custom port for accessing Whisper's services")
	flags.String(hydraAdminEndpoint, "", "Hydra Admin Enpoint")
	flags.String(logLevel, "info", "Sets the Log Level to one of seven (trace, debug, info, warn, error, fatal, panic)")
}

// InitFromViper initializes the web server builder with properties retrieved from Viper.
func (b *Builder) InitFromViper(v *viper.Viper) *Builder {
	b.Port = v.GetString(port)
	b.BaseUIPath = v.GetString(baseUIPath)
	b.HydraAdminEndpoint = v.GetString(hydraAdminEndpoint)
	b.LogLevel = v.GetString(logLevel)
	logrus.Infof("Run config: %v", b)
	return b
}
