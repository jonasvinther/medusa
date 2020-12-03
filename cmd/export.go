package cmd

import (
	"fmt"
	"medusa/pkg/vaultengine"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(exportCmd)
}

var exportCmd = &cobra.Command{
	Use:   "export [file to import]",
	Short: "Export Vault secrets as yaml",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

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

		vaultAddr, _ := cmd.Flags().GetString("vault-url")
		vaultToken, _ := cmd.Flags().GetString("vault-token")
		vaultPrefix, _ := cmd.Flags().GetString("vault-prefix")

		client := vaultengine.NewClient(vaultAddr, vaultToken, vaultPrefix)

		d, _ := client.FolderExport(path)
		y, _ := client.ConvertToYaml(d)
		// j, _ := client.ConvertToJSON(d)

		fmt.Println(y)
		// fmt.Println(j)
	},
}
