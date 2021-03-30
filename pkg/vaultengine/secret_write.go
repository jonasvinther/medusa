package vaultengine

import "fmt"

// SecretWrite is used for writing data to a Vault instance
func (client *Client) SecretWrite(path string, data map[string]interface{}) {
	infix := "/data/"

	if client.engineType == "kv1" {
		infix = ""
	}

	finalPath := client.engine + infix + path

	finalData := make(map[string]interface{})
	finalData["data"] = data

	_, err := client.vc.Logical().Write(finalPath, finalData)
	if err != nil {
		fmt.Printf("Error while writing secret. %s\n", err)
	} else {
		fmt.Printf("Secret successfully written to Vault [%s] using path [%s]\n", client.addr, path)
	}
}
