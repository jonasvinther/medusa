package vaultengine

import (
	"encoding/json"
	"fmt"
)

// SecretWrite is used for writing data to a Vault instance
func (client *Client) SecretWrite(path string, data map[string]interface{}) {
	infix := "/data/"

	if client.engineType == "kv1" {
		infix = ""
	}

	finalPath := client.engine + infix + path

	finalData := make(map[string]interface{})

	if client.engineType == "kv1" {
		finalData = data
	} else {
		finalData["data"] = data
	}

	if val, ok := data["json-key"]; ok {
		fmt.Println("Do json", val)
		b, err := json.Marshal([]byte(val.(string)))
		_, err = client.vc.Logical().WriteBytes(finalPath, b)
		if err != nil {
			fmt.Printf("Error while writing secret. %s\n", err)
		} else {
			fmt.Printf("Secret successfully written to Vault [%s] using path [%s]\n", client.addr, path)
		}
	} else {
		fmt.Printf("No json: %+v \n", data)
		_, err := client.vc.Logical().Write(finalPath, finalData)
		if err != nil {
			fmt.Printf("Error while writing secret. %s\n\n", err)
		} else {
			fmt.Printf("Secret successfully written to Vault [%s] using path [%s]\n\n", client.addr, path)
		}
	}
}
