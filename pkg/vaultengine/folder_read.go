package vaultengine

import (
	"fmt"
)

//FolderRead reads the subpaths and secrets of the provided path
func (client *Client) FolderRead(path string) ([]interface{}, error) {
	infix := "/metadata/"

	if client.engineType == "kv1" {
		infix = ""
	}

	finalPath := client.engine + infix + path

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
