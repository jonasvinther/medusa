package cmd

import (
	"medusa/pkg/importer"
	"medusa/pkg/vaultengine"

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

		vault := vaultengine.VaultEngine{
			Token:  vaultToken,
			URL:    vaultURL,
			Prefix: vaultPrefix}

		parsedYaml, _ := importer.ImportYaml(file)

		// Write the data to Vault using the Vault engine
		for path, value := range parsedYaml {
			vault.WriteSecret(path, value)
		}
	},
}
