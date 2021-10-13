# Using Docker to run Medusa

We will here show one was of using `medusa` with Docker, hiding a lot of the details and making your medusa commands shorter and faster to type.

First, define a function and export it. This could be done in your environment so that it is always available. 

```
function medusa(){
  docker run   \
    -v $(pwd):/tmp/output/  \
    --user $(id -u):$(id -u)  \
    -e VAULT_ADDR=$VAULT_ADDR \
    -e VAULT_TOKEN=$VAULT_TOKEN  \
    -e VAULT_SKIP_VERIFY=$VAULT_SKIP_VERIFY  \
    ghcr.io/jonasvinther/medusa:latest "$@"
}
export -f medusa
```

Now you can export the `VAULT_ADDR`, `VAULT_TOKEN` and `VAULT_SKIP_VERIFY`variables

```
VAULT_ADDR=https://192.168.86.41:8201
VAULT_SKIP_VERIFY=true
VAULT_TOKEN=00000000-0000-0000-0000-000000000000
```

And now you are ready to run medusa like this

```
medusa export secret/A
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

And if you want it to a file, you can use the volume

```
medusa export secret/A -o /tmp/output/backup.medusa
```

> TIP : You can change the function to use a different location inside the container that is easier to remember, instead of `/tmp/output`.