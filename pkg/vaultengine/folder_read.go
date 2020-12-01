package vaultengine

import (
	"fmt"
)

//FolderRead reads the provided path and all sub paths
// func (client *Client) FolderRead(path string) (map[string]map[string]interface{}, error) {

// 	out := make(map[string]map[string]interface{})

// 	return out, nil
// }

//FolderRead reads the provided path and all sub paths
func (client *Client) FolderRead(path string) ([]interface{}, error) {
	finalPath := client.namespace + path

	secret, err := client.vc.Logical().List(finalPath)
	if err != nil {
		return nil, err
	}

	if secret == nil {
		return nil, fmt.Errorf("no keys found using path [%s] on Vault instance [%s]", finalPath, client.addr)
	}

	keys := secret.Data["keys"].([]interface{})
	return keys, nil
}
