package vaultengine

import "log"

// SecretWrite is used for writing data to a Vault instance
func (client *Client) SecretWrite(path string, data map[string]interface{}) {
	finalPath := client.engine + "/data/" + path

	finalData := make(map[string]interface{})
	finalData["data"] = data

	_, err := client.vc.Logical().Write(finalPath, finalData)
	if err != nil {
		log.Printf("Error while writing secret. %s", err)
	} else {
		log.Printf("Secret successfully written to Vault instance on path [%s]", path)
	}
}
