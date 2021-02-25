package cmd

import (
	"fmt"
	"log"
	"medusa/pkg/encrypt"
	"medusa/pkg/vaultengine"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.PersistentFlags().StringP("format", "f", "yaml", "Specify the export format [yaml|json]")
	exportCmd.PersistentFlags().StringP("output", "o", "", "Write to file instead of stdout")
	exportCmd.PersistentFlags().BoolP("encrypt", "e", false, "Encrypt the exported Vault data")
	exportCmd.PersistentFlags().StringP("public-key", "p", "", "Location of the RSA public key")
}

var exportCmd = &cobra.Command{
	Use:   "export [vault path]",
	Short: "Export Vault secrets as yaml",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		vaultAddr, _ := cmd.Flags().GetString("address")
		vaultToken, _ := cmd.Flags().GetString("token")
		insecure, _ := cmd.Flags().GetBool("insecure")
		doEncrypt, _ := cmd.Flags().GetBool("encrypt")
		exportFormat, _ := cmd.Flags().GetString("format")
		output, _ := cmd.Flags().GetString("output")

		engine, path := vaultengine.PathSplitPrefix(path)
		client := vaultengine.NewClient(vaultAddr, vaultToken, insecure)
		client.UseEngine(engine)

		exportData, err := client.FolderExport(path)
		if err != nil {
			log.Printf("%s", err)
			return
		}

		// Convert export to json or yaml
		var data []byte
		switch exportFormat {
		case "json":
			data, err = vaultengine.ConvertToJSON(exportData)
		case "yaml":
			data, err = vaultengine.ConvertToYaml(exportData)
		default:
			fmt.Printf("Wrong format '%s' specified. Available formats are yaml and json.\n", exportFormat)
			return
		}

		if err != nil {
			fmt.Println(err)
			return
		}

		if doEncrypt {
			publicKeyPath, _ := cmd.Flags().GetString("public-key")
			encryptedKey, encryptedData := encrypt.Encrypt(publicKeyPath, output, data)

			if output == "" {
				fmt.Println(string([]byte(encryptedData)))
				fmt.Println(string(encryptedKey))
			} else {
				// Write to file
				// First encrypted data
				vaultengine.WriteToFile(output, []byte(encryptedData))
				vaultengine.AppendStringToFile(output, "\n")
				// Then encrypted AES key
				vaultengine.AppendStringToFile(output, encryptedKey)
				vaultengine.AppendStringToFile(output, "\n")
			}
		} else {
			if output == "" {
				fmt.Println(string(data))
			} else {
				vaultengine.WriteToFile(output, data)
			}
		}
	},
}
