package importer

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v2"
)

// RawYaml is a interface for arbitrary data
type RawYaml map[interface{}]interface{}

// ParsedYaml is a map for the parsed data
// The format will be: map["path/to/secret"]["secret_key"]=interface{}
type ParsedYaml map[string]map[string]interface{}

func parseYaml(rawYaml RawYaml, parsedYaml *ParsedYaml, path string) {
	for key, value := range rawYaml {
		// Handle nil values in the yaml data
		if value == nil {
			value = ""
		}

		// Check if the given object is of the same type as the RawYaml data type
		// If true - We know that we have not reached the last element of the structure yet
		if reflect.TypeOf(value).String() == reflect.TypeOf(make(RawYaml)).String() {
			tmpPath := fmt.Sprintf("%s/%s", path, key)
			parseYaml(value.(RawYaml), parsedYaml, tmpPath)
		} else {
			// Check if the key exists in the data structure
			// If it doesn't we create it
			if _, exist := (*parsedYaml)[path]; !exist {
				(*parsedYaml)[path] = make(map[string]interface{})
			}

			// Append the value to the parsed data structure using it's absolute path
			(*parsedYaml)[path][fmt.Sprintf("%v", key)] = value
		}
	}
}

// ImportYaml parses the byte stream into yaml
func ImportYaml(byteValue []byte) (parsedYaml ParsedYaml, err error) {
	// Create objects for the the two data types
	// that are needed in order to pass the yaml data
	parsedYaml = make(ParsedYaml)
	rawYaml := make(RawYaml)

	err = yaml.Unmarshal(byteValue, &rawYaml)
	if err != nil {
		return nil, err
	}

	// Parse the yaml data using a recursive loop
	parseYaml(rawYaml, &parsedYaml, "")

	return parsedYaml, nil
}
