package vaultengine

import (
	"log"
)

// SecretRead is used for reading a secret from a Vault instance
func (client *Client) SecretRead(path string) map[string]interface{} {
	infix := "/data/"

	if client.engineType == "kv1" {
		infix = "/"
	}

	finalPath := client.engine + infix + path

	secret, err := client.vc.Logical().Read(finalPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if secret == nil {
		log.Fatalf("No secret found using path [%s] on Vault instance [%s]. Medusa will terminate now.", path, client.addr)
	}

	// kv1
	if client.engineType == "kv1" {
		return secret.Data
	}

	// kv2
	m, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		// If we are using the kv2 secret engine and the current version
		// of the secret has been deleted we return nil because there are
		// no active version of the secret
		metadata := secret.Data["metadata"].(map[string]interface{})
		if metadata["deletion_time"] != "" {
			return nil
		} else {
			log.Fatalf("Error while reading secret\nPath:\t%s\nData:\t%#v\n\n", finalPath, secret.Data["data"])
		}
	}

	return m
}
