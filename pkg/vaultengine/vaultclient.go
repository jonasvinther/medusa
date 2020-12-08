package vaultengine

import (
	vault "github.com/hashicorp/vault/api"
)

// Client describes the arguments that is needed to to establish a connecting to a Vault instance
type Client struct {
	token     string
	addr      string
	namespace string
	insecure  bool
	vc        *vault.Client
}

// NewClient creates a instance of the VaultClient struct
func NewClient(addr, token, namespace string, insecure bool) *Client {

	client := &Client{
		token:     token,
		addr:      addr,
		insecure:  insecure,
		namespace: namespace}

	client.newVaultClient()

	return client
}

func (client *Client) newVaultClient() error {
	config := vault.Config{Address: client.addr}

	// Enable insecure
	config.ConfigureTLS(&vault.TLSConfig{
		Insecure: client.insecure,
	})

	vc, err := vault.NewClient(&config)
	if err != nil {
		return err
	}

	client.vc = vc

	if client.namespace != "" {
		client.vc.SetNamespace(client.namespace)
	}

	if client.token != "" {
		client.vc.SetToken(client.token)
	}

	return nil
}
