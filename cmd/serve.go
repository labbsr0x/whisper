package cmd

import (
	"github.com/abilioesteves/whisper/web"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return web.Initialize()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
