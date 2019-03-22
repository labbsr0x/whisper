package cmd

import (
	"github.com/abilioesteves/whisper/web"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server and serves the HTTP REST API",
	RunE: func(cmd *cobra.Command, args []string) error {
		return web.Initialize()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
