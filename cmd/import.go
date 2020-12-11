package cmd

import (
	"medusa/pkg/importer"
	"medusa/pkg/vaultengine"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import [vault path] [file to import]",
	Short: "Import a yaml file into a Vault instance",
	Long:  ``,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		file := args[1]
		vaultAddr, _ := cmd.Flags().GetString("address")
		vaultToken, _ := cmd.Flags().GetString("token")
		insecure, _ := cmd.Flags().GetBool("insecure")

		engine, prefix := vaultengine.PathSplitPrefix(path)
		client := vaultengine.NewClient(vaultAddr, vaultToken, insecure)
		client.UseEngine(engine)
		parsedYaml, _ := importer.ImportYaml(file)

		// Write the data to Vault using the Vault engine
		for path, value := range parsedYaml {
			path = strings.TrimPrefix(path, "/")
			path = prefix + path
			client.SecretWrite(path, value)
		}
	},
}
