# Medusa

[![GoDoc](https://godoc.org/github.com/jonasvinther/medusa?status.svg)](https://godoc.org/github.com/jonasvinther/medusa)
[![Go Report Card](https://goreportcard.com/badge/github.com/jonasvinther/medusa)](https://goreportcard.com/report/github.com/jonasvinther/medusa)
[![Build status](https://github.com/jonasvinther/medusa/workflows/Go/badge.svg)](https://github.com/jonasvinther/medusa/actions)
[![codecov](https://codecov.io/gh/jonasvinther/medusa/branch/master/graph/badge.svg)](https://codecov.io/gh/jonasvinther/medusa)

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

Now you can open Vault in a browser [https://localhost:8201](https://localhost:8201) and login with the root token (found by `echo $VAULT_TOKEN`)

## How to contribute
Please read and follow our [contributing guide](docs/CONTRIBUTING.md)
