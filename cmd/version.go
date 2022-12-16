package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "development"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Medusa",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("medusa %s", Version)

		return nil
	},
}
