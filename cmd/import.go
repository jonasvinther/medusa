package cmd

import (
	"fmt"
	"medusa/pkg/encrypt"
	"medusa/pkg/importer"
	"medusa/pkg/vaultengine"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.PersistentFlags().BoolP("decrypt", "d", false, "Decrypt the Vault data before importing")
	importCmd.PersistentFlags().StringP("private-key", "p", "", "Location of the RSA private key")
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
		doDecrypt, _ := cmd.Flags().GetBool("decrypt")
		privateKey, _ := cmd.Flags().GetString("private-key")

		engine, prefix := vaultengine.PathSplitPrefix(path)
		client := vaultengine.NewClient(vaultAddr, vaultToken, insecure)
		client.UseEngine(engine)

		var parsedYaml importer.ParsedYaml
		// var err error

		if doDecrypt {
			// Decrypt the data before parsing
			decryptedData, err := encrypt.Decrypt(privateKey, file)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Import and parse the data
			parsedYaml, _ = importer.Import([]byte(decryptedData))
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			// Read unencrypted data from file
			data, err := importer.ReadFromFile(file)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Import and parse the data
			parsedYaml, err = importer.Import(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		// Write the parsed yaml to Vault using the Vault engine
		for path, value := range parsedYaml {
			path = prefix + strings.TrimPrefix(path, "/")
			client.SecretWrite(path, value)
		}
	},
}
