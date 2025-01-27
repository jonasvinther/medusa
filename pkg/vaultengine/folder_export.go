package vaultengine

import (
	"fmt"
	"strings"
	"time"
)

// Folder defines a level of the Vault structure
type Folder map[string]interface{}

// FolderExport will export all subfolders and secrets from a specified location
func (client *Client) FolderExport(path string, sleepMillis int) (Folder, error) {
	baseFolder := make(Folder)
	subFolders := make(Folder)

	err := client.PathReader(&subFolders, path, sleepMillis)
	if err != nil {
		return nil, err
	}

	path = strings.TrimSuffix(path, "/")
	parts := strings.Split(path, "/")

	buildFolderStructure(&baseFolder, parts, subFolders)

	return baseFolder, nil
}

// buildFolderStructure creates the base tree structure
func buildFolderStructure(parentFolder *Folder, parts []string, subFolders Folder) error {
	nextPart := parts[0]
	parts = parts[1:]
	newSubFolder := make(Folder)

	if len(parts) == 0 {
		// If we are at the root level we overwrite the rootfolder with it's subfolder
		// so that we don't get empty keys in our export
		if nextPart == "" {
			*parentFolder = subFolders
		} else {
			(*parentFolder)[nextPart] = subFolders
		}

	} else {
		buildFolderStructure(&newSubFolder, parts, subFolders)
		(*parentFolder)[nextPart] = newSubFolder
	}

	return nil
}

// PathReader recursively reads the provided path and all subpaths
func (client *Client) PathReader(parentFolder *Folder, path string, sleepMillis int) error {

	folder, err := client.FolderRead(path, sleepMillis)
	if err != nil {
		return err
	}

	for _, key := range folder {
		strKey := fmt.Sprintf("%v", key)
		newPath := path + strKey

		if IsFolder(strKey) {
			subFolder := make(Folder)
			keyName := strings.Replace(strKey, "/", "", -1)

			err = client.PathReader(&subFolder, newPath, sleepMillis)
			if err != nil {
				return err
			}

			if sleepMillis > 0 {
				time.Sleep(time.Duration(sleepMillis) * time.Millisecond)
			}

			if (*parentFolder)[keyName] != nil {
				for key, elem := range (*parentFolder)[keyName].(map[string]interface{}) {
					subFolder[key] = elem
				}
			}
			(*parentFolder)[keyName] = subFolder
		} else {
			s := client.SecretRead(newPath)

			if sleepMillis > 0 {
				time.Sleep(time.Duration(sleepMillis) * time.Millisecond)
			}

			if len(s) > 0 {
				(*parentFolder)[strKey] = s
			}
		}

	}

	return nil
}
