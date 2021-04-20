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

	if client.engineType == "kv1" {
		return secret.Data
	} else {
		m, ok := secret.Data["data"].(map[string]interface{})
		if !ok {
			log.Fatalf("Error while reading secret\nPath:\t%s\nData:\t%#v\n\n", finalPath, secret.Data["data"])
		}

		return m
	}
}
