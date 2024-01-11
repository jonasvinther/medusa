package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/jonasvinther/medusa/pkg/encrypt"
	"github.com/jonasvinther/medusa/pkg/importer"
	"github.com/jonasvinther/medusa/pkg/vaultengine"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.PersistentFlags().BoolP("decrypt", "d", false, "Decrypt the Vault data before importing")
	importCmd.PersistentFlags().StringP("private-key", "p", "", "Location of the RSA private key")
	importCmd.PersistentFlags().StringP("engine-type", "m", "kv2", "Specify the secret engine type [kv1|kv2]")
}

var importCmd = &cobra.Command{
	Use:   "import [vault path] ['file' to import | '-' read from stdin]",
	Short: "Import a yaml file into a Vault instance",
	Long:  ``,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		file := args[1]
		vaultAddr, _ := cmd.Flags().GetString("address")
		vaultToken, _ := cmd.Flags().GetString("token")
		insecure, _ := cmd.Flags().GetBool("insecure")
		vaultRole, _ := cmd.Flags().GetString("role")
		kubernetes, _ := cmd.Flags().GetBool("kubernetes")
		authPath, _ := cmd.Flags().GetString("kubernetes-auth-path")
		namespace, _ := cmd.Flags().GetString("namespace")
		engineType, _ := cmd.Flags().GetString("engine-type")
		doDecrypt, _ := cmd.Flags().GetBool("decrypt")
		privateKey, _ := cmd.Flags().GetString("private-key")

		client := vaultengine.NewClient(vaultAddr, vaultToken, insecure, namespace, vaultRole, kubernetes, authPath)
		engine, prefix, err := client.MountpathSplitPrefix(path)
		if err != nil {
			fmt.Println(err)
			return err
		}

		client.UseEngine(engine)
		client.SetEngineType(engineType)

		var parsedYaml importer.ParsedYaml

		if doDecrypt {
			// Decrypt the data before parsing
			decryptedData, err := encrypt.Decrypt(privateKey, file)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// Import and parse the data
			parsedYaml, err = importer.Import([]byte(decryptedData))
			if err != nil {
				fmt.Println(err)
				return err
			}
		} else {
			var data []byte

			if file == "-" {
				// Read unencrypted data from stdin
				var inputReader io.Reader = cmd.InOrStdin()
				data, _ = ioutil.ReadAll(inputReader)
			} else {
				// Read unencrypted data from file
				data, err = importer.ReadFromFile(file)
				if err != nil {
					fmt.Println(err)
					return err
				}
			}

			// Import and parse the data
			parsedYaml, err = importer.Import(data)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}

		// Write the parsed yaml to Vault using the Vault engine
		for path, value := range parsedYaml {
			path = prefix + strings.TrimPrefix(path, "/")
			client.SecretWrite(path, value)
		}

		return nil
	},
}
