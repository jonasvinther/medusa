package importer

import (
	"testing"
)

var testYaml = []byte(`A:
B:
  C:
    D:
      Da: value 1
      Db: 
        DBa: value 1
        DBb: value 2
  E:
    Ea: value 1
    Eb: value 2
`)

// expectedYaml := map[string]map[string]interface{}{
// 	"/A/B/C/D": {
// 		"Da": "Value 1",
// 	},
// 	"/A/B/C/D/Db": {
// 		"DBa": "value 1",
// 		"DBb": "value 2",
// 	},
// 	"/A/B/E": {
// 		"Ea": "value 1",
// 		"Eb": "value 2",
// 	},
// }

func TestDoImportYaml(t *testing.T) {
	_, err := doImportYaml(testYaml)

	if err != nil {
		t.Error("Failed parsing yaml")
	}
}
