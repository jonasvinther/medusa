package cmd

import (
	"fmt"
	"medusa/pkg/importer"
	"medusa/pkg/vaultengine"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if viper.IsSet("VAULT_ADDR") {
			value := viper.Get("VAULT_ADDR").(string)
			cmd.Flags().Set("vault-url", value)
		}
		if viper.IsSet("VAULT_TOKEN") {
			value := viper.Get("VAULT_TOKEN").(string)
			cmd.Flags().Set("vault-token", value)
		}
		if viper.IsSet("VAULT_TOKEN") {
			value := viper.Get("VAULT_TOKEN").(string)
			cmd.Flags().Set("vault-token", value)
		}

		vaultURL, _ := cmd.Flags().GetString("vault-url")
		vaultToken, _ := cmd.Flags().GetString("vault-token")
		vaultPrefix, _ := cmd.Flags().GetString("vault-prefix")

		fmt.Printf("Vault token: %s\n", vaultToken)

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
