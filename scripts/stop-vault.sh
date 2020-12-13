#!/bin/bash

# Set variables
echo "Stopping Vault container"
export SCRIPT_SOURCE=$(dirname "${BASH_SOURCE[0]}")
(cd $SCRIPT_SOURCE
docker-compose -f docker-compose-tls.yml down)
