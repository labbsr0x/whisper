package cmd

import (
	"fmt"

	"github.com/abilioesteves/whisper/version"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display this build's version, build time, and git hash",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:		%s\n", version.Version)
		fmt.Printf("Git Hash:	%s\n", version.GitHash)
		fmt.Printf("Build Time:	%s\n", version.BuildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
