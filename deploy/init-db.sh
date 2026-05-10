#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER llmproxy WITH PASSWORD '$LITELLM_DB_PASSWORD';
    CREATE DATABASE litellm OWNER llmproxy;
EOSQL
