package cmd

import (
	"errors"
	"fmt"

	"github.com/jonasvinther/medusa/pkg/encrypt"
	"github.com/jonasvinther/medusa/pkg/vaultengine"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.PersistentFlags().StringP("format", "f", "yaml", "Specify the export format [yaml|json]")
	exportCmd.PersistentFlags().StringP("output", "o", "", "Write to file instead of stdout")
	exportCmd.PersistentFlags().BoolP("encrypt", "e", false, "Encrypt the exported Vault data")
	exportCmd.PersistentFlags().StringP("public-key", "p", "", "Location of the RSA public key")
	exportCmd.PersistentFlags().StringP("engine-type", "m", "kv2", "Specify the secret engine type [kv1|kv2]")
}

var exportCmd = &cobra.Command{
	Use:   "export [vault path]",
	Short: "Export Vault secrets as yaml",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		vaultAddr, _ := cmd.Flags().GetString("address")
		vaultToken, _ := cmd.Flags().GetString("token")
		vaultRole, _ := cmd.Flags().GetString("role")
		kubernetes, _ := cmd.Flags().GetBool("kubernetes")
		authPath, _ := cmd.Flags().GetString("kubernetes-auth-path")
		insecure, _ := cmd.Flags().GetBool("insecure")
		namespace, _ := cmd.Flags().GetString("namespace")
		engineType, _ := cmd.Flags().GetString("engine-type")
		doEncrypt, _ := cmd.Flags().GetBool("encrypt")
		exportFormat, _ := cmd.Flags().GetString("format")
		output, _ := cmd.Flags().GetString("output")

		client := vaultengine.NewClient(vaultAddr, vaultToken, insecure, namespace, vaultRole, kubernetes, authPath)
		engine, path, err := client.MountpathSplitPrefix(path)
		if err != nil {
			fmt.Println(err)
			return err
		}

		client.UseEngine(engine)
		client.SetEngineType(engineType)

		exportData, err := client.FolderExport(path)
		if err != nil {
			fmt.Println(err)
			return err
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
			err = errors.New("invalid export format")
		}

		if err != nil {
			fmt.Println(err)
			return err
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
				err = vaultengine.WriteToFile(output, []byte(encryptedData))
				if err != nil {
					return err
				}
				err = vaultengine.AppendStringToFile(output, "\n")
				if err != nil {
					return err
				}
				// Then encrypted AES key
				err = vaultengine.AppendStringToFile(output, encryptedKey)
				if err != nil {
					return err
				}
				err = vaultengine.AppendStringToFile(output, "\n")
				if err != nil {
					return err
				}
			}
		} else {
			if output == "" {
				fmt.Println(string(data))
			} else {
				err = vaultengine.WriteToFile(output, data)
				if err != nil {
					return err
				}
			}
		}

		return nil
	},
}
