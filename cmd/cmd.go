package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "medusa",
	Short: "Medusa is a cli tool currently for importing a json or yaml file into HashiCorp Vault.",
	Long: `Medusa is a cli tool currently for importing a json or yaml file into HashiCorp Vault.
Created by by Jonas Vinther & Henrik Høegh.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Apply the viper config value to the flag when the flag is not set and viper has a value
		address, _ := cmd.Flags().GetString("address")
		if viper.IsSet("VAULT_ADDR") && address == "" {
			value := viper.Get("VAULT_ADDR").(string)
			cmd.Flags().Set("address", value)
		}

		token, _ := cmd.Flags().GetString("token")
		if viper.IsSet("VAULT_TOKEN") && token == "" {
			value := viper.Get("VAULT_TOKEN").(string)
			cmd.Flags().Set("token", value)
		}

		return nil
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("address", "a", "", "Address of the Vault server")
	rootCmd.PersistentFlags().StringP("token", "t", "", "Vault authentication token")
	rootCmd.PersistentFlags().BoolP("insecure", "k", false, "Allow insecure server connections when using SSL")

	// SetConfigFile explicitly defines the path, name and extension of the config file.
	// Viper will use this and not check any of the config paths.
	// .env - It will search for the .env file in the current directory
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.medusa")
	viper.SetConfigName("config.yaml")

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		// log.Fatalf("Error while reading config file %s", err)
	}

}
