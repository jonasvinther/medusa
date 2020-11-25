#!/bin/bash

# Set variables
echo "Stopping Vault container"
export SCRIPT_SOURCE=$(dirname "${BASH_SOURCE[0]}")
(cd $SCRIPT_SOURCE
source .env
docker-compose down)
