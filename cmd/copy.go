package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"
	"gopkg.in/yaml.v3"
	"github.com/spf13/cobra"
	"github.com/jonasvinther/medusa/pkg/vaultengine"
	"github.com/jonasvinther/medusa/pkg/importer"
)

// Initialize the "copy" command and its flags
func init() {
	rootCmd.AddCommand(copyCmd)
	copyCmd.PersistentFlags().StringP("engine-type", "m", "kv2", "Specify the secret engine type [kv1|kv2]")
}

// Define the "copy" command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy Vault secret from one path to another",
	Long:  `The "copy" command allows users to export secrets from a source path in Vault and copy them to a target path.`,
	Args:  cobra.MinimumNArgs(2), // Require at least two arguments: source path and target path
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse command-line arguments and flags
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

		client := vaultengine.NewClient(vaultAddr, vaultToken, insecure, namespace, vaultRole, kubernetes, authPath)
		engine, sourcePath, err := client.MountpathSplitPrefix(sourcePath)
		if err != nil {
			fmt.Println("Error during source path split:", err)
			return err
		}
		client.UseEngine(engine)
		client.SetEngineType(engineType)

		// Export secrets from the source path
		exportData, err := client.FolderExport(sourcePath)
		if err != nil {
			fmt.Println("Error during export:", err)
			return err
		}

		// Check if exported data is empty
		if len(exportData) == 0 {
			return fmt.Errorf("No data found at source path %s", sourcePath)
		}

		// Write exported secrets to a temporary file in YAML format
		tempFileName := "/tmp/exported_secret.yaml"
		data, err := vaultengine.ConvertToYaml(exportData)
		if err != nil {
			fmt.Println("Error during YAML conversion:", err)
			return err
		}

		err = ioutil.WriteFile(tempFileName, data, 0644)
		if err != nil {
			fmt.Println("Error writing temporary file:", err)
			return err
		}

		// Extract specific data from the YAML file based on the source path
		sourcePath = strings.TrimSuffix(sourcePath, "/")
		err = extractYamlData(tempFileName, sourcePath)
		if err != nil {
			fmt.Println("Error extracting YAML data:", err)
			return err
		}

		// Read the modified YAML file
		fileData, err := ioutil.ReadFile(tempFileName)
		if err != nil {
			fmt.Println("Error reading modified YAML file:", err)
			return err
		}

		// Import modified data into the target path
		parsedYaml, err := importer.Import(fileData)
		if err != nil {
			fmt.Println("Error importing YAML data:", err)
			return err
		}

		// Write the imported data to the target path in Vault
		for subPath, value := range parsedYaml {
			fullPath := targetPath + subPath
			client.SecretWrite(fullPath, value)
		}

		fmt.Printf("Secrets from %s have been successfully copied to %s\n", sourcePath, targetPath)
		return nil
	},
}

// Function to extract specific YAML data based on a given path
func extractYamlData(inputFile, path string) error {
	// Read the YAML file
	fileContent, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("Error reading file: %v", err)
	}

	var data map[string]interface{}
	// Parse the YAML content into a Go data structure
	if err := yaml.Unmarshal(fileContent, &data); err != nil {
		return fmt.Errorf("Error parsing YAML file: %v", err)
	}

	// Split the path into keys
	pathKeys := strings.Split(path, "/")

	// Navigate through the data structure using the keys
	for _, key := range pathKeys {
		if value, exists := data[key]; exists {
			if mapData, ok := value.(map[string]interface{}); ok {
				data = mapData
			} else {
				data = map[string]interface{}{key: value}
				break
			}
		} else {
			return fmt.Errorf("Path '%s' does not exist in the YAML file", path)
		}
	}

	// Save the extracted data into a new YAML file
	outputFile := "/tmp/exported_secret.yaml"
	outputData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("Error converting data to YAML: %v", err)
	}

	if err := ioutil.WriteFile(outputFile, outputData, 0644); err != nil {
		return fmt.Errorf("Error writing output file: %v", err)
	}

	fmt.Printf("Extracted data has been saved to '%s'.\n", outputFile)
	return nil
}
