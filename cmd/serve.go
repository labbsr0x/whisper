package cmd

import (
	"fmt"

	"github.com/abilioesteves/whisper/web"
	"github.com/abilioesteves/whisper/web/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the HTTP REST APIs server",
	RunE: func(cmd *cobra.Command, args []string) error {
		webBuilder := new(config.WebBuilder).InitFromViper(viper.GetViper())
		server := new(web.Server).InitFromWebBuilder(webBuilder)
		err := server.Run()
		if err != nil {
			return fmt.Errorf("An error occurred while setting up the Whisper Web Server: %v", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	config.AddFlags(serveCmd.Flags())
	// db.AddFlags(serveCmd.Flags()) // TODO

	err := viper.GetViper().BindPFlags(serveCmd.Flags())
	if err != nil {
		panic(err)
	}
}
