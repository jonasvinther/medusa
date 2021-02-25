package importer

import "io/ioutil"

// ReadFromFile reads the content of a file
func ReadFromFile(file string) (data []byte, err error) {
	byteValue, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return byteValue, nil
}
