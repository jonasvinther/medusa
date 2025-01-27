package vaultengine

import (
	"fmt"
	"log"

	"github.com/hashicorp/vault/api"
)

// CollectPaths will retrieve all paths to secrets defined under the given path
func (client *Client) CollectPaths(path string, sleepMillis int) ([]string, error) {
	var secretPaths []string
	folder, err := client.FolderRead(path, sleepMillis)
	if err != nil {
		return nil, err
	}

	for _, key := range folder {
		strKey := fmt.Sprintf("%v", key)
		newPath := path + strKey
		newPath = CleanupPath(newPath)

		if IsFolder(strKey) {
			t, err := client.CollectPaths(newPath, sleepMillis)
			secretPaths = append(secretPaths, t...)
			if err != nil {
				return nil, err
			}
		} else {
			// client.SecretDelete(newPath)
			secretPaths = append(secretPaths, newPath)
		}
	}

	return secretPaths, nil
}

// SecretDelete deletet a Vault secret by given path
func (client *Client) SecretDelete(path string) (*api.Secret, error) {
	infix := "/metadata"

	if client.engineType == "kv1" {
		infix = "/"
	}

	finalPath := client.engine + infix + path

	result, err := client.vc.Logical().Delete(finalPath)
	if err != nil {
		log.Fatalf("Unable to delete secret: %s", err)
	}

	return result, nil
}
