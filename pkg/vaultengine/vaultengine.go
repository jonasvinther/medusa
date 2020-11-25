package vaultengine

import (
	"log"

	"github.com/hashicorp/vault/api"
)

// VaultEngine describes the arguments that is needed to to establish a connecting to a Vault instance
type VaultEngine struct {
	Token  string
	URL    string
	Prefix string
}

// WriteSecret is used for writing data to a Vault instance
func (engine VaultEngine) WriteSecret(path string, data map[string]interface{}) {

	config := api.Config{Address: engine.URL}
	config.ConfigureTLS(&api.TLSConfig{
		Insecure: true,
	})

	client, err := api.NewClient(&config)
	if err != nil {
		log.Fatalf("Error while creating client. %s", err)
	}

	client.SetToken(engine.Token)

	finalPath := engine.Prefix + path

	finalData := make(map[string]interface{})
	finalData["data"] = data

	_, err = client.Logical().Write(finalPath, finalData)
	if err != nil {
		log.Printf("Error while writing secret. %s", err)
	} else {
		log.Printf("Secret succesfully written to Vault instance on path [%s]", path)
	}
}

// GetSecret is used for reading a secret from a Vault instance
func (engine VaultEngine) GetSecret(path string) map[string]interface{} {
	client, err := api.NewClient(&api.Config{Address: engine.URL})
	if err != nil {
		log.Fatalf("%v", err)
	}

	client.SetToken(engine.Token)

	finalPath := engine.Prefix + path

	secret, err := client.Logical().Read(finalPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if secret == nil {
		log.Fatalf("No secret found using path [%s] on Vault instance [%s]. Pandora will terminate now.", path, engine.URL)
	}

	m, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		log.Fatalf("%T %#v", secret.Data["data"], secret.Data["data"])
	}

	return m
}
