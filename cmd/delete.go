package cmd

import (
	"fmt"

	"github.com/jonasvinther/medusa/pkg/vaultengine"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.PersistentFlags().StringP("engine-type", "m", "kv2", "Specify the secret engine type [kv1|kv2]")
	deleteCmd.PersistentFlags().BoolP("auto-approve", "y", false, "Skip interactive approval of plan before deletion")
}

var deleteCmd = &cobra.Command{
	Use:   "delete [vault path] [flags]",
	Short: "Recursively delete all secrets below the given path",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		vaultAddr, _ := cmd.Flags().GetString("address")
		vaultToken, _ := cmd.Flags().GetString("token")
		insecure, _ := cmd.Flags().GetBool("insecure")
		vaultRole, _ := cmd.Flags().GetString("role")
		kubernetes, _ := cmd.Flags().GetBool("kubernetes")
		authPath, _ := cmd.Flags().GetString("kubernetes-auth-path")
		namespace, _ := cmd.Flags().GetString("namespace")
		engineType, _ := cmd.Flags().GetString("engine-type")
		isApproved, _ := cmd.Flags().GetBool("auto-approve")

		// Setup Vault client
		client := vaultengine.NewClient(vaultAddr, vaultToken, insecure, namespace, vaultRole, kubernetes, authPath)
		engine, path, err := client.MountpathSplitPrefix(path)
		if err != nil {
			fmt.Println(err)
			return err
		}

		client.UseEngine(engine)
		client.SetEngineType(engineType)

		// Recursive delete
		secretPaths, err := client.CollectPaths(path)
		if err != nil {
			return err
		}

		// Print a list of all the secrets that will be deleted
		for _, key := range secretPaths {
			fmt.Printf("Deleting secret [%s%s]\n", engine, key)
		}

		// Prompt for confirmation
		if !isApproved {
			prompt := promptui.Prompt{
				Label:     fmt.Sprintf("Do you want to delete the %d secrets listed above? Only 'y' will be accepted to approve.", len(secretPaths)),
				IsConfirm: true,
			}

			result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Aborting. No secrets got deleted\n")
			}

			if result == "y" {
				isApproved = true
			}
		}

		// Perform deletion of the secrets
		if isApproved {
			for _, key := range secretPaths {
				client.SecretDelete(key)
			}
			fmt.Printf("The secrets has now been deleted\n")
		}

		return nil
	},
}
