package vaultengine

import (
	"fmt"
	"strings"
)

// Folder defines a level of the Vault structure
type Folder map[string]interface{}

// FolderExport will export all subfolders and secrets from a specified location
func (client *Client) FolderExport(path string) (Folder, error) {
	baseFolder := make(Folder)
	subFolders := make(Folder)

	err := client.pathReader(&subFolders, path)
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

//pathReader recursively reads the provided path and all subpaths
func (client *Client) pathReader(parentFolder *Folder, path string) error {
	folder, err := client.FolderRead(path)
	if err != nil {
		return err
	}

	for _, key := range folder {
		strKey := fmt.Sprintf("%v", key)
		newPath := path + strKey

		if IsFolder(strKey) {
			subFolder := make(Folder)
			keyName := strings.Replace(strKey, "/", "", -1)

			err = client.pathReader(&subFolder, newPath)
			if err != nil {
				return err
			}

			if (*parentFolder)[keyName] != nil {
				for key, elem := range (*parentFolder)[keyName].(map[string]interface{}) {
					subFolder[key] = elem
			    }
            }
			(*parentFolder)[keyName] = subFolder
		} else {
			s := client.SecretRead(newPath)
			(*parentFolder)[strKey] = s
		}
	}

	return nil
}
