package exporter

// RawYaml is a interface for arbitrary data
type RawYaml map[interface{}]interface{}

// ExportYaml will export the secrets of a Vault secret engine into yaml
func ExportYaml() (RawYaml, error) {
	exportedYaml := make(RawYaml)

	return exportedYaml, nil
}
