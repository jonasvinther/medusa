package cmd

import (
	"fmt"
	"log"
	"medusa/pkg/vaultengine"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.PersistentFlags().StringP("format", "f", "yaml", "Specify the export format [yaml|json]")
	exportCmd.PersistentFlags().StringP("output", "o", "", "Write to file instead of stdout")
}

var exportCmd = &cobra.Command{
	Use:   "export [file to import]",
	Short: "Export Vault secrets as yaml",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		vaultAddr, _ := cmd.Flags().GetString("address")
		vaultToken, _ := cmd.Flags().GetString("token")
		insecure, _ := cmd.Flags().GetBool("insecure")
		vaultPrefix, _ := cmd.Flags().GetString("vault-prefix")
		exportFormat, _ := cmd.Flags().GetString("format")
		output, _ := cmd.Flags().GetString("output")

		client := vaultengine.NewClient(vaultAddr, vaultToken, vaultPrefix, insecure)
		d, err := client.FolderExport(path)

		if err != nil {
			log.Printf("%s", err)
			return
		}

		switch exportFormat {
		case "json":
			data, _ := client.ConvertToJSON(d)

			if output == "" {
				fmt.Println(string(data))
			} else {
				client.WriteToFile(output, data)
			}
		case "yaml":
			data, _ := client.ConvertToYaml(d)

			if output == "" {
				fmt.Println(string(data))
			} else {
				client.WriteToFile(output, data)
			}
		default:
			log.Printf("Wrong format specified")
		}
	},
}
