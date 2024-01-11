package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "medusa",
	Short: "Medusa is a cli tool currently for importing a json or yaml file into HashiCorp Vault.",
	Long: `Medusa is a cli tool currently for importing a json or yaml file into HashiCorp Vault.
Created by Jonas Vinther & Henrik Høegh.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Apply the viper config value to the flag when the flag is not set and viper has a value
		address, _ := cmd.Flags().GetString("address")
		if viper.IsSet("VAULT_ADDR") && address == "" {
			value := viper.Get("VAULT_ADDR").(string)
			err := cmd.Flags().Set("address", value)
			if err != nil {
				return err
			}
		}

		token, _ := cmd.Flags().GetString("token")
		if viper.IsSet("VAULT_TOKEN") && token == "" {
			value := viper.Get("VAULT_TOKEN").(string)
			err := cmd.Flags().Set("token", value)
			if err != nil {
				return err
			}
		}

		role, _ := cmd.Flags().GetString("role")
		if viper.IsSet("VAULT_ROLE") && role == "" {
			value := viper.Get("VAULT_ROLE").(string)
			err := cmd.Flags().Set("role", value)
			if err != nil {
				return err
			}
		}

		kubernetes, _ := cmd.Flags().GetBool("kubernetes")
		if viper.IsSet("KUBERNETES") && kubernetes == false {
			value := viper.GetBool("KUBERNETES")
			err := cmd.Flags().Set("kubernetes", strconv.FormatBool(value))
			if err != nil {
				return err
			}
		}

		insecure, _ := cmd.Flags().GetBool("insecure")
		if viper.IsSet("VAULT_SKIP_VERIFY") && insecure == false {
			value := viper.GetBool("VAULT_SKIP_VERIFY")
			err := cmd.Flags().Set("insecure", strconv.FormatBool(value))
			if err != nil {
				return err
			}
		}

		namespace, _ := cmd.Flags().GetString("namespace")
		if viper.IsSet("VAULT_NAMESPACE") && namespace == "" {
			value := viper.Get("VAULT_NAMESPACE").(string)
			err := cmd.Flags().Set("namespace", value)
			if err != nil {
				return err
			}
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
	rootCmd.PersistentFlags().StringP("role", "r", "", "Vault role for Kubernetes JWT authentication")
	rootCmd.PersistentFlags().BoolP("kubernetes", "", false, "Authenticate using the Kubernetes JWT token")
	rootCmd.PersistentFlags().StringP("kubernetes-auth-path", "", "", "Authentication mount point within Vault for Kubernetes")
	rootCmd.PersistentFlags().BoolP("insecure", "k", false, "Allow insecure server connections when using SSL")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace within the Vault server (Enterprise only)")

	// AutomaticEnv makes Viper load environment variables
	viper.AutomaticEnv()

	// Explicitly defines the path, name and type of the config file.
	// Viper will use this and not check any of the config paths.
	// It will search for the "config" file in the ~/.medusa
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.medusa")
	viper.SetConfigName("config")

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		// log.Fatalf("Error while reading config file %s", err)
	}

}
