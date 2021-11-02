package vaultengine

import (
  "errors"
  "strings"

	vault "github.com/hashicorp/vault/api"
)

// Client describes the arguments that is needed to to establish a connecting to a Vault instance
type Client struct {
	token      string
	addr       string
	namespace  string
	engine     string
	engineType string
	insecure   bool
	vc         *vault.Client
}

// NewClient creates a instance of the VaultClient struct
func NewClient(addr, token string, insecure bool, namespace string) *Client {
	client := &Client{
		token:     token,
		addr:      addr,
		insecure:  insecure,
		namespace: namespace}

	client.newVaultClient()

	return client
}

// UseEngine defines which engine the Vault client will use
func (client *Client) UseEngine(engine string) {
	client.engine = engine
}

func (client *Client) MountpathSplitPrefix(path string) (string, string, error) {
  // Split Engine mountpath from path

  r := client.vc.NewRequest("GET", "/v1/sys/internal/ui/mounts/"+path)
  resp, err := client.vc.RawRequest(r)
  if resp != nil {
    defer resp.Body.Close()
  }
  if err != nil {
    // any 404 indicates k/v v1
    if resp != nil && resp.StatusCode == 404 {
      return "", "path", nil
    }
    return "", "", err
  }

  secret, err := vault.ParseSecret(resp.Body)
  if err != nil {
    return "", "", err
  }
  if secret == nil {
    return "", "", errors.New("nil response from pre-flight request")
  }
  var mountPath string
  if mountPathRaw, ok := secret.Data["path"]; ok {
    mountPath = mountPathRaw.(string)
  }

  mountPath = strings.TrimSuffix(mountPath, "/")
  suffix := strings.Replace(path, mountPath, "", 1)
  suffix = EnsureFolder(strings.TrimPrefix(suffix, "/"))

  return mountPath, suffix, nil
}

// SetEngineType defines which vault secret engine type that is being used
func (client *Client) SetEngineType(engineType string) {
	client.engineType = engineType
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
