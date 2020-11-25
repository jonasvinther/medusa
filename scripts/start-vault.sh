#!/bin/bash

export VAULT_VOLUME=/tmp/vault/data/vault-volume
export VAULT_ADD=""
export VAULT_ROOT_TOKEN="00000000-0000-0000-0000-000000000000"
export HOST_IP=$(ip route get 1 | sed -n 's/^.*src \([0-9.]*\) .*$/\1/p')

# Generate self signed certificates for Vault to use
openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
    -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=vault.local" \
    -keyout vault.local.key  -out vault.local.crt

# Cleanup the temp folder
if [ -d "$VAULT_VOLUME" ] 
then
    rm -fr $VAULT_VOLUME/*
else
    mkdir -p $VAULT_VOLUME
fi

# Move certificates to temp folder
mv vault.local.key /tmp/vault/data/vault-volume/
mv vault.local.crt /tmp/vault/data/vault-volume/

# Change permission on Vault volume folder
chmod 700 $VAULT_VOLUME/*

# Setting UID to 100 for the certificates for Vault to use
docker run -d -u root -v $VAULT_VOLUME:/tmp/vault:rw alpine:latest chown -R 100:1000 ls /tmp/vault

# Start the Vault container
docker-compose up -d

# Wait for Vault to startup
sleep 5

# Extract the container ip for Vault and output to terminal
VAULT_IP=$(docker inspect vault --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}')
echo "Vault container ip is : $VAULT_IP"

# Run a `Vault status` to check if it is alive
docker run --network=container:vault --cap-add IPC_LOCK -e VAULT_ADDR=https://$VAULT_IP:8201 -e VAULT_SKIP_VERIFY=true --rm vault:latest vault status

# Echo help on how to interact with Vault from the same host
echo "
To run vault commands, use the following docker command:
 - For connecting to Vault on the same Docker host:
     docker run --network=container:vault --cap-add IPC_LOCK -e VAULT_ADDR=https://$VAULT_IP:8201 -e VAULT_SKIP_VERIFY=true --rm vault:latest vault status 
 - For connecting to Vault from another host :
     docker run --network=container:vault --cap-add IPC_LOCK -e VAULT_ADDR=https://$HOST_IP:8201 -e VAULT_SKIP_VERIFY=true --rm vault:latest vault status 
"

# Generate .env file
echo "export VAULT_ADDR=https://$HOST_IP:8201
export VAULT_SKIP_VERIFY=true
export VAULT_TOKEN=$VAULT_ROOT_TOKEN" > .env
