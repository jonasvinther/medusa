package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"
	"github.com/spf13/cobra"
	"github.com/jonasvinther/medusa/pkg/vaultengine"
	"github.com/jonasvinther/medusa/pkg/importer"
	"github.com/manifoldco/promptui"

)

func init() {
	rootCmd.AddCommand(moveCmd)
	moveCmd.PersistentFlags().BoolP("auto-approve", "y", false, "Skip interactive approval of plan before deletion")
	moveCmd.PersistentFlags().StringP("engine-type", "m", "kv2", "Specify the secret engine type [kv1|kv2]")
}

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move Vault secret from one path to another",
	Long:  ``,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourcePath := args[0]
		targetPath := args[1]
		vaultAddr, _ := cmd.Flags().GetString("address")
		vaultToken, _ := cmd.Flags().GetString("token")
		vaultRole, _ := cmd.Flags().GetString("role")
		kubernetes, _ := cmd.Flags().GetBool("kubernetes")
		authPath, _ := cmd.Flags().GetString("kubernetes-auth-path")
		insecure, _ := cmd.Flags().GetBool("insecure")
		namespace, _ := cmd.Flags().GetString("namespace")
		engineType, _ := cmd.Flags().GetString("engine-type")
		isApproved, _ := cmd.Flags().GetBool("auto-approve")


		client := vaultengine.NewClient(vaultAddr, vaultToken, insecure, namespace, vaultRole, kubernetes, authPath)
		engine, sourcePath, err := client.MountpathSplitPrefix(sourcePath)

		if err != nil {
			fmt.Println("error splitting source path:", err)
			return err
		}
		client.UseEngine(engine)
		client.SetEngineType(engineType)


		exportData, err := client.FolderExport(sourcePath)
		if err != nil {
			fmt.Println("error exporting data:", err)
			return err
		}


		if len(exportData) == 0 {
			return fmt.Errorf("no data found in source path %s", sourcePath)
		}


		tempFileName := "/tmp/exported_secret.yaml"
		data, err := vaultengine.ConvertToYaml(exportData)
		if err != nil {
			fmt.Println("error converting data to YAML:", err)
			return err
		}

		err = ioutil.WriteFile(tempFileName, data, 0644)
		if err != nil {
			fmt.Println("error writing temporary file:", err)
			return err
		}

		sourcePath_edited := strings.TrimSuffix(sourcePath, "/")
		err = extractYamlData(tempFileName, sourcePath_edited)
		if err != nil {
			fmt.Println("error editing YAML data:", err)
			return err
		}


		fileData, err := ioutil.ReadFile(tempFileName)
		if err != nil {
			fmt.Println("error reading modified file:", err)
			return err
		}

		parsedYaml, err := importer.Import(fileData)
		if err != nil {
			fmt.Println("error parsing YAML data:", err)
			return err
		}

		for subPath, value := range parsedYaml {
			fullPath := targetPath + subPath		
			client.SecretWrite(fullPath, value)
		}

		secretPaths, err := client.CollectPaths(sourcePath)
		if err != nil {
			return err
		}


		for _, key := range secretPaths {
			fmt.Printf("Deleting secret [%s%s]\n", engine, key)
		}

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
		if isApproved {
			for _, key := range secretPaths {
				client.SecretDelete(key)
			}
			fmt.Printf("The secrets has now been deleted\n")
		}

		return nil
		fmt.Printf("Secrets moved from %s to %s successfully.\n", sourcePath, targetPath)
		return nil
	},
}

