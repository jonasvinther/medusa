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
	Use:   "export [vault path]",
	Short: "Export Vault secrets as yaml",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		vaultAddr, _ := cmd.Flags().GetString("address")
		vaultToken, _ := cmd.Flags().GetString("token")
		insecure, _ := cmd.Flags().GetBool("insecure")
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

		switch exportFormat {
		case "json":
			data, _ := vaultengine.ConvertToJSON(exportData)

			if output == "" {
				fmt.Println(string(data))
			} else {
				vaultengine.WriteToFile(output, data)
			}
		case "yaml":
			data, _ := vaultengine.ConvertToYaml(exportData)

			if output == "" {
				fmt.Println(string(data))
			} else {
				vaultengine.WriteToFile(output, data)
			}
		default:
			log.Printf("Wrong format specified")
		}
	},
}
