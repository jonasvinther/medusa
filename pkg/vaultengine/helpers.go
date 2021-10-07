package vaultengine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

// IsFolder returns true if the path is a folder
// Folders are suffixed with "/" in Vault
func IsFolder(p string) bool {
	return strings.HasSuffix(p, "/")
}

// ConvertToYaml will return the Folder object as yaml
func ConvertToYaml(folder Folder) ([]byte, error) {
	yaml, err := yaml.Marshal(folder)
	if err != nil {
		return nil, err
	}
	return yaml, nil
}

// ConvertToJSON will return the Folder object as yaml
func ConvertToJSON(folder Folder) ([]byte, error) {
	json, err := json.MarshalIndent(folder, "", "\t")
	// json, err := json.Marshal(folder)
	if err != nil {
		return nil, err
	}
	return json, nil
}

// WriteToFile will create a file and store the provieded data in it
func WriteToFile(file string, data []byte) error {
	// file, _ := json.MarshalIndent(data, "", " ")

	err := ioutil.WriteFile(file, data, 0644)

	return err
}

/**
 * Append string to the end of file
 *
 * path: the path of the file
 * text: the string to be appended. If you want to append text line to file,
 *       put a newline '\n' at the of the string.
 */
func AppendStringToFile(path, text string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}

// PathSplitPrefix will split the first part of a string into it's own variable
// and return it together with the rest of the path
func PathSplitPrefix(path string) (string, string) {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	parts := strings.Split(path, "/")
	prefix := parts[0]
	suffix := strings.Join(parts[1:], "/")
	suffix = EnsureFolder(suffix)
	return prefix, suffix
}

// PathJoin combines multiple paths into one.
func PathJoin(p ...string) string {
	if strings.HasSuffix(p[len(p)-1], "/") {
		return strings.TrimPrefix(path.Join(p...)+"/", "/")
	}
	return strings.TrimPrefix(path.Join(p...), "/")
}

// CleanupPath cleans up unwanted characters so that the path looks nice and clean
func CleanupPath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}
	return strings.ReplaceAll(path, "//", "/")
}

// EnsureFolder ensures a path is a folder (adds a trailing "/").
func EnsureFolder(p string) string {
	return PathJoin(p, "/")
}
