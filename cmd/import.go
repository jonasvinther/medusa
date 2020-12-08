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
		vaultAddr, _ := cmd.Flags().GetString("address")
		vaultToken, _ := cmd.Flags().GetString("token")
		insecure, _ := cmd.Flags().GetBool("insecure")
		vaultPrefix, _ := cmd.Flags().GetString("vault-prefix")

		client := vaultengine.NewClient(vaultAddr, vaultToken, vaultPrefix, insecure)

		parsedYaml, _ := importer.ImportYaml(file)

		// Write the data to Vault using the Vault engine
		for path, value := range parsedYaml {
			client.SecretWrite(path, value)
		}
	},
}
