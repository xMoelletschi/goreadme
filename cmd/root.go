package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goreadme",
	Short: "Generate useful boilerplate and docs for Go projects",
	Long: `goreadme is a CLI tool to help generate structured documentation,
boilerplate files, and config templates for Go repositories.

Examples:
  goreadme go             Generate a README.md from Go documentation

Use "--help" on any subcommand for more options.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {}
