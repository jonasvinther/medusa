package cmd

import (
	"fmt"

	"github.com/jonasvinther/medusa/pkg/encrypt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.PersistentFlags().StringP("private-key", "p", "", "Location of the RSA private key")
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt [file path] [flags]",
	Short: "Decrypt an encrypted Vault output file into plaintext in stdout",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]
		privateKey, _ := cmd.Flags().GetString("private-key")

		decryptedData, err := encrypt.Decrypt(privateKey, file)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Printf("%s", decryptedData)

		return nil
	},
}
