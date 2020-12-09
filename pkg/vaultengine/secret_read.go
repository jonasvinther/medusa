package vaultengine

import "log"

// SecretRead is used for reading a secret from a Vault instance
func (client *Client) SecretRead(path string) map[string]interface{} {
	finalPath := client.engine + "/data/" + path

	secret, err := client.vc.Logical().Read(finalPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if secret == nil {
		log.Fatalf("No secret found using path [%s] on Vault instance [%s]. Medusa will terminate now.", path, client.addr)
	}

	m, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		log.Fatalf("%T %#v", secret.Data["data"], secret.Data["data"])
	}

	return m
}
