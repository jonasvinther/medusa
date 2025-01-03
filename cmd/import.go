package cmd
 
import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
	Use:   "import [vault path] ['file1' 'file2' ... or '-' to read from stdin]",
	Short: "Import yaml/json files or folder of secrets into a Vault instance",
	Long:  "",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		inputPaths := args[1:]
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

		for _, inputPath := range inputPaths {
			if isDir(inputPath) {
				// Process directory recursively
				err := processDirectory(inputPath, func(filePath string) error {
					return processFile(filePath, cmd, client, prefix, doDecrypt, privateKey)
				})
				if err != nil {
					fmt.Printf("Error importing directory %s: %v\n", inputPath, err)
				}
			} else {
				// Process single file or stdin
				err := processFile(inputPath, cmd, client, prefix, doDecrypt, privateKey)
				if err != nil {
					fmt.Printf("Error importing file %s: %v\n", inputPath, err)
				}
			}
		}

		return nil
	},
}

// Check if the path is a directory
func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Process a directory recursively
func processDirectory(path string, fileHandler func(string) error) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories themselves
		if info.IsDir() {
			return nil
		}
		// Handle individual files
		return fileHandler(filePath)
	})
}

// Process an individual file
func processFile(file string, cmd *cobra.Command, client *vaultengine.Client, prefix string, doDecrypt bool, privateKey string) error {
	var data []byte
	var err error

	if file == "-" {
		// Read data from stdin
		var inputReader io.Reader = cmd.InOrStdin()
		data, err = ioutil.ReadAll(inputReader)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %v", err)
		}
	} else {
		// Read data from file
		data, err = importer.ReadFromFile(file)
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", file, err)
		}
	}

	// Decrypt if required
	if doDecrypt {
		decryptedData, err := encrypt.Decrypt(privateKey, file)
		if err != nil {
			return fmt.Errorf("error decrypting file %s: %v", file, err)
		}
		data = []byte(decryptedData)
	}

	// Import and parse the data
	parsedYaml, err := importer.Import(data)
	if err != nil {
		return fmt.Errorf("error importing file %s: %v", file, err)
	}

	// Write data to Vault
	for subPath, value := range parsedYaml {
		fullPath := prefix + strings.TrimPrefix(subPath, "/")
		client.SecretWrite(fullPath, value)
	}

	return nil
}
