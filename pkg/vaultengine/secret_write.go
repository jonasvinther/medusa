package vaultengine

import (
	"fmt"
)

// SecretWrite is used for writing data to a Vault instance
func (client *Client) SecretWrite(path string, data map[string]interface{}) {
	infix := "/data/"

	if client.engineType == "kv1" {
		infix = "/"
	}

	finalPath := client.engine + infix + path

	finalData := make(map[string]interface{})

	if client.engineType == "kv1" {
		finalData = data
	} else {
		finalData["data"] = data
	}

	// If the data object has the json-object key
	// it means that the secret is not in the default
	// key/value format.
	if jsonVal, ok := data["json-object"]; ok {
		var jsonString string

		// The kv2 engine needs the data wrapped in a "data" key
		if client.engineType == "kv2" {
			jsonString = fmt.Sprintf("{\"data\":%s}", jsonVal)
		} else {
			jsonString = jsonVal.(string)
		}

		_, err := client.vc.Logical().WriteBytes(finalPath, []byte(jsonString))
		if err != nil {
			fmt.Printf("Error while writing secret. %s\n", err)
		} else {
			fmt.Printf("Secret successfully written to Vault [%s] using path [%s]\n", client.addr, path)
		}
	} else {
		_, err := client.vc.Logical().Write(finalPath, finalData)
		if err != nil {
			fmt.Printf("Error while writing secret. %s\n", err)
		} else {
			fmt.Printf("Secret successfully written to Vault [%s] using path [%s]\n", client.addr, path)
		}
	}
}
