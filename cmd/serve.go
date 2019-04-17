package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/labbsr0x/whisper-client/client"
	whisperClientConfig "github.com/labbsr0x/whisper-client/config"

	"github.com/labbsr0x/whisper/web"
	"github.com/labbsr0x/whisper/web/config"
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
		self := new(client.WhisperClient).InitFromHydraClient(server.HydraClient)
		t, err := self.CheckCredentials()

		if err == nil {
			token := self.GetTokenAsJSONStr(t)
			logrus.Debugf("Initial Token: '%v'", token)
			os.Setenv(string(whisperClientConfig.WhisperTokenEnvKey), token) // now the token can be referenced by a middleware or any other execution point

			err = server.Run()
			if err != nil {
				return fmt.Errorf("An error occurred while setting up the Whisper Web Server: %v", err)
			}
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
