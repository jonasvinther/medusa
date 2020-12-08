package vaultengine

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// IsFolder returns true if the path is a folder
// Folders are suffixed with "/" in Vault
func IsFolder(p string) bool {
	return strings.HasSuffix(p, "/")
}

// ConvertToYaml will return the Folder object as yaml
func (client *Client) ConvertToYaml(folder Folder) ([]byte, error) {
	yaml, err := yaml.Marshal(folder)
	if err != nil {
		return nil, err
	}
	return yaml, nil
}

// ConvertToJSON will return the Folder object as yaml
func (client *Client) ConvertToJSON(folder Folder) ([]byte, error) {
	json, err := json.MarshalIndent(folder, "", "\t")
	// json, err := json.Marshal(folder)
	if err != nil {
		return nil, err
	}
	return json, nil
}

// WriteToFile will create a file and store the provieded data in it
func (client *Client) WriteToFile(filename string, data []byte) error {
	// file, _ := json.MarshalIndent(data, "", " ")

	err := ioutil.WriteFile(filename, data, 0644)

	return err
}
