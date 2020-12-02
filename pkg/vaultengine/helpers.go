package vaultengine

import "strings"

// IsFolder returns true if the path is a folder
// Folders are suffixed with "/" in Vault
func IsFolder(p string) bool {
	return strings.HasSuffix(p, "/")
}
