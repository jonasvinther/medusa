package cmd

import (
	"medusa/pkg/importer"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import [file to import]",
	Short: "Import a yaml file into a Vault instance",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		vaultURL, _ := cmd.Flags().GetString("vault-url")
		vaultToken, _ := cmd.Flags().GetString("vault-token")
		vaultPrefix, _ := cmd.Flags().GetString("vault-prefix")

		vault := importer.VaultEngine{
			Token:  vaultToken,
			URL:    vaultURL,
			Prefix: vaultPrefix}

		vault.ImportYaml(file)
	},
}
