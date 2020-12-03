package vaultengine

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v2"
)

// IsFolder returns true if the path is a folder
// Folders are suffixed with "/" in Vault
func IsFolder(p string) bool {
	return strings.HasSuffix(p, "/")
}

// ConvertToYaml will return the Folder object as yaml
func (client *Client) ConvertToYaml(folder Folder) (string, error) {
	yaml, err := yaml.Marshal(folder)
	if err != nil {
		return "", err
	}
	return string(yaml), nil
}

// ConvertToJSON will return the Folder object as yaml
func (client *Client) ConvertToJSON(folder Folder) (string, error) {
	json, err := json.Marshal(folder)
	if err != nil {
		return "", err
	}
	return string(json), nil
}
