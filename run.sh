#!/bin/bash

docker-compose up

# Checking if db container was successfully started
if [ $? -eq 0 ]; then
    echo "DB running completed successfully. Running migrations."
else
    echo "Error during running db. Exiting script."
    exit 1
fi

docker exec -i avito-db-1 psql -U postgres -d postgres < ./schema/000001_init.up.sql

# Applying migrations
if [ $? -eq 0 ]; then
    echo "Migrations applied successfully to the database."
else
    echo "Error while running migrations."
fi