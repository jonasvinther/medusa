# Medusa

[![GoDoc](https://godoc.org/github.com/jonasvinther/medusa?status.svg)](https://godoc.org/github.com/jonasvinther/medusa)
[![Go Report Card](https://goreportcard.com/badge/github.com/jonasvinther/medusa)](https://goreportcard.com/report/github.com/jonasvinther/medusa)
[![Build status](https://github.com/jonasvinther/medusa/workflows/Go/badge.svg)](https://github.com/jonasvinther/medusa/actions)
[![codecov](https://codecov.io/gh/jonasvinther/medusa/branch/main/graph/badge.svg)](https://codecov.io/gh/jonasvinther/medusa)

## About
Medusa is a cli tool currently for importing a json or yaml file into HashiCorp Vault.

## How to use
Go learn how to use the various commands, check out [quick-guide](docs/quick-guide.md)

## Medusa help
To test out `medusa` on your laptop
```
Medusa is a cli tool currently for importing a json or yaml file into HashiCorp Vault.
Created by Jonas Vinther & Henrik Høegh.

Usage:
  medusa [command]

Available Commands:
  export      Export Vault secrets as yaml
  help        Help about any command
  import      Import a yaml file into a Vault instance

Flags:
  -a, --address string   Address of the Vault server
  -h, --help             help for medusa
  -k, --insecure         Allow insecure server connections when using SSL
  -t, --token string     Vault authentication token

Use "medusa [command] --help" for more information about a command.
``` 

## How to contribute
Please read and follow our [contributing guide](docs/CONTRIBUTING.md)
