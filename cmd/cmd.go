package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "medusa",
	Short: "Medusa is a cli tool currently for importing a json or yaml file into HashiCorp Vault.",
	Long: `Medusa is a cli tool currently for importing a json or yaml file into HashiCorp Vault.
Created by by Jonas Vinther & Henrik Høegh.`,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("vault-url", "v", "undefined", "Vault url")
	rootCmd.PersistentFlags().StringP("vault-token", "t", "undefined", "Vault token")
	rootCmd.PersistentFlags().StringP("vault-prefix", "p", "", "Prefix of the Vault path")

	// SetConfigFile explicitly defines the path, name and extension of the config file.
	// Viper will use this and not check any of the config paths.
	// .env - It will search for the .env file in the current directory
	viper.SetConfigFile("scripts/.env")

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	viper.BindPFlag("listen.bind", rootCmd.Flags().Lookup("VAULT_ADDR"))

	// viper.Get() returns an empty interface{}
	// to get the underlying type of the key,
	// we have to do the type assertion, we know the underlying value is string
	// if we type assert to other type it will throw an error
	// value, _ := viper.Get("VAULT_ADDR").(string)

	// fmt.Println(value)

}
