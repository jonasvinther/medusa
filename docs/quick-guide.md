# Medusa quick guide

## Setting up Medusa
### Config file
It's possible to create a config file for Medusa to read in your homefolder `~/.medusa/config.yaml` that looks like this
```
VAULT_ADDR: https://192.168.86.41:8201
VAULT_SKIP_VERIFY: true
VAULT_TOKEN: 00000000-0000-0000-0000-000000000000
```
If you haven't set any environment variables, or given any parameters, this file will tell Medusa where to connect, the token to use and to `VAULT_SKIP_VERIFY` should be enabled or not.

### Environment variables
It's also possible configure Medusa via environment variables by setting them like this:
```
export VAULT_ADDR=https://192.168.86.41:8201
export VAULT_SKIP_VERIFY=true
export VAULT_TOKEN=00000000-0000-0000-0000-000000000000
```

### Parameters
> Get help with `./medusa -h`
You can configure Medusa in the commands you run like this :
```
  -a, --address string   Address of the Vault server
  -k, --insecure         Allow insecure server connections when using SSL
  -t, --token string     Vault authentication token
```

Use them like this:
```
/medusa import secret ./test/data/import-example-1.yaml -a="https://0.0.0.0:8201" -t="00000000-0000-0000-0000-000000000000" --insecure
./medusa export secret/A -a="https://0.0.0.0:8201" -t="00000000-0000-0000-0000-000000000000" --format="json" --insecure
```

## Importing data
> Get help with `./medusa import -h`
Medusa import will take a [vault path] with [flags]
Example:
```
./medusa import secret ./test/data/import-example-1.yaml -a="https://0.0.0.0:8201" -t="00000000-0000-0000-0000-000000000000" --insecure
2020/12/11 13:23:59 Secret successfully written to Vault instance on path [/A/B/E]
2020/12/11 13:23:59 Secret successfully written to Vault instance on path [/A/Xa/Z]
2020/12/11 13:23:59 Secret successfully written to Vault instance on path [/A/F/G]
2020/12/11 13:23:59 Secret successfully written to Vault instance on path [/A/B/C/D]
2020/12/11 13:23:59 Secret successfully written to Vault instance on path [/A/B/C/D/Db]



./medusa import secret/folder ./test/data/import-example-1.yaml -a="https://0.0.0.0:8201" -t="00000000-0000-0000-0000-000000000000" --insecure
2020/12/11 13:25:03 Secret successfully written to Vault instance on path [folder/A/F/G]
2020/12/11 13:25:03 Secret successfully written to Vault instance on path [folder/A/B/C/D]
2020/12/11 13:25:03 Secret successfully written to Vault instance on path [folder/A/B/C/D/Db]
2020/12/11 13:25:03 Secret successfully written to Vault instance on path [folder/A/B/E]
2020/12/11 13:25:03 Secret successfully written to Vault instance on path [folder/A/Xa/Z]

```

## Exporting data
> Get help with `./medusa export -h` and yaml is the default output format
Medusa import will take a [vault path] with [flags]
Example:

```
./medusa export secret -a="https://0.0.0.0:8201" -t="00000000-0000-0000-0000-000000000000" --format="yaml" --insecure
"":
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
  folder:
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

If you want to get the output in a Json format, simply add `format="json"` like this
```
./medusa export secret -a="https://0.0.0.0:8201" -t="00000000-0000-0000-0000-000000000000" --format="json" --insecure
{
  "": {
    "A": {
      "B": {
        "C": {
          "D": {
            "Db": {
              "DBa": "value 1",
              "DBb": "value 2"
            }
          }
        },
        "E": {
          "Ea": "value 1",
          "Eb": "value 2"
        }
      },
      "F": {
        "G": {
          "Ga": "value1"
        }
      },
      "Xa": {
        "Z": {
          "Za": "value 1",
          "Zb": "value 2"
        }
      }
    },
    "folder": {
      "A": {
        "B": {
          "C": {
            "D": {
              "Db": {
                "DBa": "value 1",
                "DBb": "value 2"
              }
            }
          },
          "E": {
            "Ea": "value 1",
            "Eb": "value 2"
          }
        },
        "F": {
          "G": {
            "Ga": "value1"
          }
        },
        "Xa": {
          "Z": {
            "Za": "value 1",
            "Zb": "value 2"
          }
        }
      }
    }
  }
}
```
