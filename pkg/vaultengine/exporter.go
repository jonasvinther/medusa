package vaultengine

import (
	"fmt"
	"strings"
)

// ParsedYaml is a map for the parsed data
// The format will be: map["path/to/secret"]["secret_key"]=interface{}
type ParsedYaml map[string]interface{}

// ExportYaml will export the secrets of a Vault secret engine into yaml
func (client *Client) ExportYaml(path string) (ParsedYaml, error) {
	exportedYaml := make(ParsedYaml)
	t := make(ParsedYaml)

	client.ExportSecrets(&t, path)

	p := strings.Replace(path, "/", "", -1)
	exportedYaml[p] = t

	return exportedYaml, nil
}

//ExportSecrets recursively reads the provided path and all subpaths
func (client *Client) ExportSecrets(exportedYaml *ParsedYaml, path string) {

	secret, err := client.FolderRead(path)

	if err != nil {
		fmt.Println(err)
	} else {
		for _, key := range secret {
			strKey := fmt.Sprintf("%v", key)

			if IsFolder(strKey) {
				newPath := path + strKey
				tmp := make(ParsedYaml)
				tmpKey := strings.Replace(strKey, "/", "", -1)
				client.ExportSecrets(&tmp, newPath)
				(*exportedYaml)[tmpKey] = tmp
			} else {
				newPath := path + strKey
				s := client.SecretRead(newPath)
				(*exportedYaml)[strKey] = s
			}
		}
	}
}
