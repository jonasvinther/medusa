# Medusa

## About
Medusa is a cli tool currently for importing a json or yaml file into HashiCorp Vault.

## How to use
To test out `medusa` on your laptop
```
# Build the binary
go build

# Start a local Vault instance in Docker
./scripts/start-vault.sh

# Setup environment variables to point to Vault
source scripts/.env

# Import a test yaml file to Vault
./medusa import ./test/data/import-example-1.yaml -p="secret/data" -v="$VAULT_ADDR" -t="$VAULT_TOKEN"
``` 

Now you can open Vault in a browser `https://localhost:8201` and login with the root token (found by `echo $VAULT_TOKEN`)

## How to contribute
Create an issue or pull request
