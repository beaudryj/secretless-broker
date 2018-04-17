#!/bin/bash -ex

source ./_conjur.sh

admin_api_key=$(docker-compose exec conjur conjurctl role retrieve-key dev:user:admin | tr -d '\r')
myapp_api_key=$(conjur_cli "$admin_api_key" host rotate_api_key -h myapp | tr -d '\r')

export CONJUR_AUTHN_API_KEY="$myapp_api_key"

docker-compose up -d myapp_secretless

export DB_HOST=/var/lib/postgresql

docker-compose run --rm \
  --entrypoint ./makedb.sh \
  myapp

docker-compose up --no-deps -d myapp