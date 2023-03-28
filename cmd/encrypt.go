package cmd

import (
	"fmt"
	"os"

	"github.com/jonasvinther/medusa/pkg/encrypt"
	"github.com/jonasvinther/medusa/pkg/vaultengine"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.PersistentFlags().StringP("output", "o", "", "Write to file instead of stdout")
	encryptCmd.PersistentFlags().StringP("public-key", "p", "", "Location of the RSA public key")
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt [file path] [flags]",
	Short: "Encrypt a Vault export file onto stdout or to an output file",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]
		publicKey, _ := cmd.Flags().GetString("public-key")
		output, _ := cmd.Flags().GetString("output")

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Println(err)
			return err
		}

		encryptedKey, encryptedData := encrypt.Encrypt(publicKey, output, data)

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

		return nil
	},
}
