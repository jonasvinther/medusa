![medusa logo](https://raw.githubusercontent.com/jonasvinther/medusa/main/assets/logo/medusa-icon-240.png#gh-light-mode-only)
![medusa logo](https://raw.githubusercontent.com/jonasvinther/medusa/main/assets/logo/medusa-icon-240-white.png#gh-dark-mode-only)


# Medusa

[![GoDoc](https://godoc.org/github.com/jonasvinther/medusa?status.svg)](https://godoc.org/github.com/jonasvinther/medusa)
[![Go Report Card](https://goreportcard.com/badge/github.com/jonasvinther/medusa)](https://goreportcard.com/report/github.com/jonasvinther/medusa)
[![Build status](https://github.com/jonasvinther/medusa/workflows/Go/badge.svg)](https://github.com/jonasvinther/medusa/actions)
[![codecov](https://codecov.io/gh/jonasvinther/medusa/branch/main/graph/badge.svg)](https://codecov.io/gh/jonasvinther/medusa)
[![CodeQL](https://github.com/jonasvinther/medusa/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/jonasvinther/medusa/actions/workflows/codeql-analysis.yml)

## Table of Contents

- [About](#about)
- [Supported HashiCorp Vault versions](#supported-hashicorp-vault-versions)
- [How to use](#how-to-use)
  * [Setting up Medusa](#setting-up-medusa)
  * [Importing secrets](#importing-secrets)
  * [Improved import secrets](#improved-import-secrets)
  * [Exporting secrets](#exporting-secrets)
  * [Deleting secrets](#deleting-secrets)
  * [Copying secrets](#copying-secrets)
  * [Moving secrets](#moving-secrets)
  * [Decrypt secrets](#decrypt-secrets)
  * [Kubernetes examples](docs/examples/kubernetes/cronjob/)
  * [Docker examples](docs/examples/docker/)
- [Secure secret management outside Vault](#secure-secret-management-outside-vault)
- [Help](#help)
- [How to contribute](docs/CONTRIBUTING.md)

## About
Medusa is a cli tool currently for importing and exporting a json or yaml file into HashiCorp Vault.  
Medusa currently supports kv1 and kv2 Vault secret engines.

## Supported HashiCorp Vault versions
The minimum required HashiCorp Vault version that is supported by Medusa is Vault version 0.10.0. Medusa has not been tested with earlier verisons than Vault version 0.10.0.

## How to use
In this section you can read about how to configure and use Medusa.  
You can also watch the [Medusa 101 introduction video](https://youtu.be/ynoe3fs_YHg) to get a quick introduction on how to use Medusa for importing and exporting secrets in HashiCorp Vault.

### Setting up Medusa
#### Config file
It's possible to create a config file for Medusa to read in your homefolder `~/.medusa/config.yaml` that looks like this
```
VAULT_ADDR: https://192.168.86.41:8201
VAULT_SKIP_VERIFY: true
VAULT_TOKEN: 00000000-0000-0000-0000-000000000000
```
If you haven't set any environment variables, or given any parameters, this file will tell Medusa where to connect, the token to use and to `VAULT_SKIP_VERIFY` should be enabled or not.

#### Environment variables
It's also possible configure Medusa via environment variables by setting them like this:
```
export VAULT_ADDR=https://192.168.86.41:8201
export VAULT_SKIP_VERIFY=true
export VAULT_TOKEN=00000000-0000-0000-0000-000000000000
```

#### Parameters
> Get help with `./medusa -h`
You can configure Medusa in the commands you run like this :
```
  -a, --address string   Address of the Vault server
  -k, --insecure         Allow insecure server connections when using SSL
  -t, --token string     Vault authentication token
```

Use them like this:
```
./medusa import secret ./test/data/import-example-1.yaml --address="https://0.0.0.0:8201" --token="00000000-0000-0000-0000-000000000000" --insecure
./medusa export secret/A --address="https://0.0.0.0:8201" --token="00000000-0000-0000-0000-000000000000" --format="json" --insecure
```

### Importing secrets
> Get help with `./medusa import -h`

Import a yaml file into a Vault instance

Usage:
  medusa import [vault path] ['file' to import | '-' read from stdin] [flags]

```
  Flags:
  -d, --decrypt              Decrypt the Vault data before importing
  -m, --engine-type string   Specify the secret engine type [kv1|kv2] (default "kv2")
  -h, --help                 help for import
  -p, --private-key string   Location of the RSA private key

  Global Flags:
  -a, --address string       Address of the Vault server
  -k, --insecure             Allow insecure server connections when using SSL
  -n, --namespace string     Namespace within the Vault server (Enterprise only)
  -t, --token string         Vault authentication token
```

Example:
```bash
# Read from file
./medusa import secret ./test/data/import-example-1.yaml -a="https://0.0.0.0:8201" -t="00000000-0000-0000-0000-000000000000" --insecure
Secret successfully written to Vault instance on path [/A/B/E]
Secret successfully written to Vault instance on path [/A/Xa/Z]
Secret successfully written to Vault instance on path [/A/F/G]
Secret successfully written to Vault instance on path [/A/B/C/D]
Secret successfully written to Vault instance on path [/A/B/C/D/Db]

# Read from file
./medusa import secret/folder ./test/data/import-example-1.yaml -a="https://0.0.0.0:8201" -t="00000000-0000-0000-0000-000000000000" --insecure
Secret successfully written to Vault instance on path [folder/A/F/G]
Secret successfully written to Vault instance on path [folder/A/B/C/D]
Secret successfully written to Vault instance on path [folder/A/B/C/D/Db]
Secret successfully written to Vault instance on path [folder/A/B/E]
Secret successfully written to Vault instance on path [folder/A/Xa/Z]

# Read from stdin
cat ./test/data/import-example-1.yaml | ./medusa import secret -
Secret successfully written to Vault [https://localhost:8201] using path [/A/F/G]
Secret successfully written to Vault [https://localhost:8201] using path [/A/B/C/D]
Secret successfully written to Vault [https://localhost:8201] using path [/A/B/C/D/Db]
Secret successfully written to Vault [https://localhost:8201] using path [/A/B/E]
Secret successfully written to Vault [https://localhost:8201] using path [/A/Xa/Z]
```

### Exporting secrets
> Get help with `./medusa export -h` and yaml is the default output format
Medusa import will take a [vault path] with [flags]

```
  Flags:
      --display-keys-only    Display only keys of secrets but not their values
  -e, --encrypt              Encrypt the exported Vault data
  -m, --engine-type string   Specify the secret engine type [kv1|kv2] (default "kv2")
  -f, --format string        Specify the export format [yaml|json] (default "yaml")
  -h, --help                 help for export
  -o, --output string        Write to file instead of stdout
  -p, --public-key string    Location of the RSA public key

  Global Flags:
  -a, --address string       Address of the Vault server
  -k, --insecure             Allow insecure server connections when using SSL
  -n, --namespace string     Namespace within the Vault server (Enterprise only)
  -t, --token string         Vault authentication token
```

Example:

```
./medusa export secret --address="https://0.0.0.0:8201" --token="00000000-0000-0000-0000-000000000000" --format="yaml" --insecure
A:
  B:
    C:
      D:
        Db:
          DBa: value 1
          DBb: value 2
    E:
      Ea: value 1
      Eb: value 2
  F:
    G:
      Ga: value1
  Xa:
    Z:
      Za: value 1
      Zb: value 2
```

### Deleting secrets
> Get help with `./medusa delete -h`
Medusa delete will take a [vault path] with [flags]

```
  Flags:
  -y, --auto-approve         Skip interactive approval of plan before deletion
  -m, --engine-type string   Specify the secret engine type [kv1|kv2] (default "kv2")
  -h, --help                 help for import

  Global Flags:
  -a, --address string       Address of the Vault server
  -k, --insecure             Allow insecure server connections when using SSL
  -n, --namespace string     Namespace within the Vault server (Enterprise only)
  -t, --token string         Vault authentication token
```

Example:
```
./medusa delete secret/production --address="https://0.0.0.0:8201" --token="00000000-0000-0000-0000-000000000000" --insecure
Deleting secret [secret/production/users/cart/database]
Deleting secret [secret/production/users/cart/database/users/readuser]
Deleting secret [secret/production/users/cart/database/users/writeuser]
Deleting secret [secret/production/users/user/database]
Deleting secret [secret/production/users/user/database/users/readuser]
? Do you want to delete the 25 secrets listed above? Only 'y' will be accepted to approve.? [y/N] y
The secrets has now been deleted


./medusa delete secret/staging --address="https://0.0.0.0:8201" --token="00000000-0000-0000-0000-000000000000" --insecure --auto-approve
Deleting secret [secret/staging/users/cart/database]
Deleting secret [secret/staging/users/cart/database/users/readuser]
Deleting secret [secret/staging/users/cart/database/users/writeuser]
Deleting secret [secret/staging/users/user/database]
Deleting secret [secret/staging/users/user/database/users/readuser]
The secrets has now been deleted
```

### Decrypt secrets
> Get help with `./medusa decrypt -h`
Medusa decrypt will take a [FILE path] with [flags]

```
  Flags:
  -p, --private-key string   Location of the RSA private key
```

Example:
```
# Write to stdout
./medusa decrypt encrypted-export.txt --private-key private-key.pem
env:
  dev:
    nomad:
      token: secret-token
  production:
    nomad:
      token: secret-other-token


# Write to file
./medusa decrypt encrypted-export.txt --private-key private-key.pem > plaintext-export.yaml
```

### Encrypt secrets
> Get help with `./medusa encrypt -h`
Medusa encrypt will take a [FILE path] with [flags]

```
  Flags:
  -o, --output string       Write to file instead of stdout
  -p, --public-key string   Location of the RSA public key
```

Example:
```
# Write to stdout
./medusa encrypt plaintext-export.txt --public-key public-key.pem
&lt;Encrypted data&gt;

# Write to file
./medusa encrypt plaintext-export.txt --public-key public-key.pem --output encrypted-export.txt.b64
```

## Secure secret management outside Vault
Medusa will help you securely manage your secrets outside Vault.
This could for instance be as a backup of your Vault data or while your secrets are being transported between Vault instances.  
Medusa uses a hybrid encryption solution in order to keep your secrets safe.  

### Key generation
When exporting your Vault secrets using Medusa, the secrets are encrypted using the AES symmetric encryption algorithm. The 256-bit AES encryption key is randomly generated by Medusa every time the export command is being called.  
Then the AES key is encrypted by the provided RSA public key and then stored together with the encrypted secrets.  
This ensures that both the exported secrets and AES enctyption key can be transfered safely between Vault instances.  
The exported secrets and AES enctyption key can only be decrypted by a person who is in possession of the RSA private key.

The RSA key-pair can be generated by the following two commands:
``` bash
# Generate private key
openssl genrsa -out private-key.pem 4096

# Generate public key
openssl rsa -in private-key.pem -pubout -out public-key.pem
```

### Exporting and encrypting Vault secrets
Encrypting your Vault export is easy using Medusa. Simply add the following two flags to your command:

```
-e, --encrypt bool       Encrypt the exported Vault data [true/false]
-p, --public-key string  Location of the RSA public key
```

Use them like this:
``` bash
./medusa export kv --address="https://my-vault-server.com" --token="00000000-0000-0000-0000-000000000000" --insecure --encrypt="true" --public-key="public-key.pem" --output="encrypted-vault-secrets.txt"
```

### Importing and decrypting Vault secrets
Decrypting and importing your encrypted Vault export can be done by adding the following two flags to your command:

```
-d, --decrypt bool        Decrypt the Vault data before importing [true/false]
-p, --private-key string  Location of the RSA private key
```

Use them like this:
``` bash
./medusa import kv encrypted-vault-secrets.txt --address="https://my-vault-server.com" --token="00000000-0000-0000-0000-000000000000" --insecure --decrypt="true" --private-key="private-key.pem"
```

## Help
To test out `medusa` on your laptop
```
Medusa is a cli tool currently for importing a json or yaml file into HashiCorp Vault.
Created by Jonas Vinther & Henrik Høegh.

Usage:
  medusa [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  decrypt     Decrypt an encrypted Vault output file into plaintext in stdout
  delete      Recursively delete all secrets below the given path
  encrypt     Encrypt a Vault export file onto stdout or to an output file
  export      Export Vault secrets as yaml
  help        Help about any command
  import      Import a yaml file into a Vault instance
  version     Print the version number of Medusa

Flags:
  -a, --address string                Address of the Vault server
  -h, --help                          help for medusa
  -k, --insecure                      Allow insecure server connections when using SSL
      --kubernetes                    Authenticate using the Kubernetes JWT token
      --kubernetes-auth-path string   Authentication mount point within Vault for Kubernetes
  -n, --namespace string              Namespace within the Vault server (Enterprise only)
  -r, --role string                   Vault role for Kubernetes JWT authentication
  -t, --token string                  Vault authentication token

Use "medusa [command] --help" for more information about a command.
``` 
### Copying secrets
> Get help with `./medusa copy -h`

The "copy" command allows users to export secrets from a source path in Vault and copy them to a target path.
```
Usage:
  medusa copy [flags]

Flags:
  -m, --engine-type string   Specify the secret engine type [kv1|kv2] (default "kv2")
  -h, --help                 help for copy

Global Flags:
  -a, --address string                Address of the Vault server
  -k, --insecure                      Allow insecure server connections when using SSL
      --kubernetes                    Authenticate using the Kubernetes JWT token
      --kubernetes-auth-path string   Authentication mount point within Vault for Kubernetes
  -n, --namespace string              Namespace within the Vault server (Enterprise only)
  -r, --role string                   Vault role for Kubernetes JWT authentication
  -t, --token string                  Vault authentication token
```

Example:
```bash
# Read from file
./medusa copy secret/A/B A/C/B -a="http://vault:8201" -t="00000000-0000-0000-0000-000000000000" --insecure
Extracted data has been saved to '/tmp/exported_secret.yaml'.
Secret successfully written to Vault [http://vault:8201] using path [A/C/B/C/D]
Secret successfully written to Vault [http://vault:8201] using path [A/C/B/C/D/Db]
Secret successfully written to Vault [http://vault:8201] using path [A/C/B/E]
Secrets from A/B have been successfully copied to A/C/B
```
The secrets on path secret/A/B are copied to secret/A/C/B

### Moving secrets
> Get help with `./medusa move -h`

Move Vault secret from one path to another
```
Usage:
  medusa move [flags]

Flags:
  -y, --auto-approve         Skip interactive approval of plan before deletion
  -m, --engine-type string   Specify the secret engine type [kv1|kv2] (default "kv2")
  -h, --help                 help for move

Global Flags:
  -a, --address string                Address of the Vault server
  -k, --insecure                      Allow insecure server connections when using SSL
      --kubernetes                    Authenticate using the Kubernetes JWT token
      --kubernetes-auth-path string   Authentication mount point within Vault for Kubernetes
  -n, --namespace string              Namespace within the Vault server (Enterprise only)
  -r, --role string                   Vault role for Kubernetes JWT authentication
  -t, --token string                  Vault authentication token
```

Example:
```bash
# Read from file
./medusa move secret/A/B A/C/B -a="http://vault:8201" -t="00000000-0000-0000-0000-000000000000" --insecure

Extracted data has been saved to '/tmp/exported_secret.yaml'.
Secret successfully written to Vault [http://vault:8201] using path [A/C/B/C/D/Db]
Secret successfully written to Vault [http://vault:8201] using path [A/C/B/E]
Secret successfully written to Vault [http://vault:8201] using path [A/C/B/C/D]
Deleting secret [secret/A/B/C/D]
Deleting secret [secret/A/B/C/D/Db]
Deleting secret [secret/A/B/E]
Do you want to delete the 3 secrets listed above? Only 'y' will be accepted to approve.: y
The secrets has now been deleted
```
The secrets on path secret/A/B are moved to secret/A/C/B

### Improved import secrets
Now it is possible to import one or multiple files with a single command. It is also possible to import an entire folder containing secret files. Medusa will first ensure that the files in the folder comply with the Vault secret format, then proceed to import them to the path specified in the command.

> Get help with `./medusa import -h`

Import yaml/json files or folder of secrets into a Vault instance
```bash
Usage:
  medusa import [vault path] ['file1' 'file2' ... or '-' to read from stdin] [flags]

Flags:
  -d, --decrypt              Decrypt the Vault data before importing
  -m, --engine-type string   Specify the secret engine type [kv1|kv2] (default "kv2")
  -h, --help                 help for import
  -p, --private-key string   Location of the RSA private key

Global Flags:
  -a, --address string                Address of the Vault server
  -k, --insecure                      Allow insecure server connections when using SSL
      --kubernetes                    Authenticate using the Kubernetes JWT token
      --kubernetes-auth-path string   Authentication mount point within Vault for Kubernetes
  -n, --namespace string              Namespace within the Vault server (Enterprise only)
  -r, --role string                   Vault role for Kubernetes JWT authentication
  -t, --token string                  Vault authentication token
```
Example for multiple files:
```bash
# Read from file
./medusa import secret test/data/example-dev-prod-env.yaml test/data/example-different-datatypes.yaml -a="http://vault:8201" -t="00000000-0000-0000-0000-000000000000" --insecure

Secret successfully written to Vault [http://vault:8201] using path [/production/env/platform/nats]
Secret successfully written to Vault [http://vault:8201] using path [/production/users/cart/database/users/writeuser]
Secret successfully written to Vault [http://vault:8201] using path [/production/users/user/database/users/readuser]
Secret successfully written to Vault [http://vault:8201] using path [/production/users/user/database]
Secret successfully written to Vault [http://vault:8201] using path [/production/env/platform/vault]
Secret successfully written to Vault [http://vault:8201] using path [/dev/users/root]
Secret successfully written to Vault [http://vault:8201] using path [/dev/users/user/database/users/readuser]
Secret successfully written to Vault [http://vault:8201] using path [/production/users/cart/database/users/readuser]
...
```
Example for secrets folder:
```bash
./medusa import secret test/data -a="http://vault:8201" -t="00000000-0000-0000-0000-000000000000" --insecure

Secret successfully written to Vault [http://vault:8201] using path [/production/users/root]
Secret successfully written to Vault [http://vault:8201] using path [/dev/users/ui/database/users/readuser]
Secret successfully written to Vault [http://vault:8201] using path [/production/env/platform/vault]
Secret successfully written to Vault [http://vault:8201] using path [/dev/users/user/database]
Secret successfully written to Vault [http://vault:8201] using path [/dev/users/cart/database/users/writeuser]
Secret successfully written to Vault [http://vault:8201] using path [/dev/users/root]
Secret successfully written to Vault [http://vault:8201] using path [/production/users/cart/database/users/writeuser]
Secret successfully written to Vault [http://vault:8201] using path [/production/users/user/database/users/readuser]
...
```