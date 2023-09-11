#!/bin/bash

# Set variables
export SCRIPT_SOURCE=$(dirname "${BASH_SOURCE[0]}")
export VAULT_VOLUME=/tmp/vault/data/vault-volume
export VAULT_TOKEN="00000000-0000-0000-0000-000000000000"

if ! [ -x "$(command -v ip)" ]; then
  echo "The command ip is not found, reverting to localhost"
  export HOST_IP="localhost"
else
    export HOST_IP=$(ip route get 1 | sed -n 's/^.*src \([0-9.]*\) .*$/\1/p')
fi

export VAULT_ADDR="https://$HOST_IP:8201"

# Cleanup the temp folder
if [ -d "$VAULT_VOLUME" ] 
then
    echo "removing old certificate folder"
    # Setting UID to 100 for the certificates for Vault to use
    docker run -d -u root -v /:/tmp/vault:rw alpine:latest rm -rf /tmp/vault/tmp/vault
fi
    echo "Creating folder for certificates"
    mkdir -p $VAULT_VOLUME

# Generate self signed certificates for Vault to use
openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
    -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=vault.local" \
    -keyout $VAULT_VOLUME/vault.local.key  -out $VAULT_VOLUME/vault.local.crt

# Change permission on Vault volume folder
chmod 700 $VAULT_VOLUME/*

# Setting UID to 100 for the certificates for Vault to use
docker run -d -u root -v $VAULT_VOLUME:/tmp/vault:rw alpine:latest chown -R 100:1000 ls /tmp/vault

# Start the Vault container
(cd $SCRIPT_SOURCE
docker-compose -f docker-compose-tls.yml up -d)

# Wait for Vault to startup
sleep 5

# Extract the container ip for Vault and output to terminal
VAULT_IP=$(docker inspect vault --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}')
echo "Vault container ip is : $VAULT_IP"

# Run a `Vault status` to check if it is alive
docker run --network=container:vault --cap-add IPC_LOCK -e VAULT_ADDR=https://$VAULT_IP:8201 -e VAULT_SKIP_VERIFY=true --rm vault:1.13.3 vault status

# Echo help on how to interact with Vault from the same host
echo "
To run vault commands, use the following docker command:
 - For connecting to Vault on the same Docker host:
     docker run --network=container:vault --cap-add IPC_LOCK -e VAULT_ADDR=https://$VAULT_IP:8201 -e VAULT_SKIP_VERIFY=true --rm vault:1.13.3 vault status 
 - For connecting to Vault from another host :
     docker run --network=container:vault --cap-add IPC_LOCK -e VAULT_ADDR=https://$HOST_IP:8201 -e VAULT_SKIP_VERIFY=true --rm vault:1.13.3 vault status 
"

# Generate .env file
# Cleanup the temp folder
if [ ! -d "~/.medusa" ]; then
  mkdir -p ~/.medusa
fi
echo "Generating Medusa config file"

echo "VAULT_ADDR: $VAULT_ADDR
VAULT_SKIP_VERIFY: true
VAULT_TOKEN: $VAULT_TOKEN" > ~/.medusa/config.yaml
