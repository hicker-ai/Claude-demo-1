package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "usermanager",
	Short: "User management system with LDAP server",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/config.yaml", "config file path")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
