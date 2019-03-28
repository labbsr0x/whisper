package web

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	baseUIPath    = "base-ui-path"
	port          = "port"
	hydraEndpoint = "hydra-endpoint"
)

// AddFlags adds flags for Builder.
func AddFlags(flags *pflag.FlagSet) {
	flags.String(baseUIPath, "", "Base path where the 'static' folder will be found with all the UI files")
	flags.String(port, "7070", "Custom port for accessing Whisper's services")
	flags.String(hydraEndpoint, "", "Hydra enpoint")
}

// InitFromViper initializes the web server builder with properties retrieved from Viper.
func (b *Builder) InitFromViper(v *viper.Viper) *Builder {
	b.Port = v.GetString(port)
	b.BaseUIPath = v.GetString(baseUIPath)
	b.HydraEndpoint = v.GetString(hydraEndpoint)
	return b
}
