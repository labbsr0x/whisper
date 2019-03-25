package cmd

import (
	"fmt"

	"github.com/abilioesteves/whisper/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the HTTP REST APIs server",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverBuilder := new(web.Builder).InitFromViper(viper.GetViper())
		server, err := serverBuilder.New()
		if err != nil {
			return fmt.Errorf("An error occurred while setting up the Whisper Web Server: %v", err)
		}
		return server.Initialize()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	web.AddFlags(serveCmd.Flags())
	// db.AddFlags(serveCmd.Flags()) // TODO

	err := viper.GetViper().BindPFlags(serveCmd.Flags())
	if err != nil {
		panic(err)
	}
}
