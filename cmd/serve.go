package cmd

import (
	"github.com/labbsr0x/whisper/mail"
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
		mailChannel := make(chan mail.Mail)

		builder := new(config.WebBuilder).Init(viper.GetViper(), mailChannel)
		
		defer builder.DB.Close()

		mailHandler := new(mail.DefaultHandler).Init(builder.MailUser, builder.MailPassword, builder.MailHost, builder.MailPort, mailChannel)
		mailHandler.Run()

		server := new(web.Server).InitFromWebBuilder(builder)

		_, err := server.Self.CheckCredentials()
		if err != nil {
			return err
		}

		err = server.Run()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	config.AddFlags(serveCmd.Flags())

	err := viper.GetViper().BindPFlags(serveCmd.Flags())
	if err != nil {
		panic(err)
	}
}
