#!/bin/bash

docker-compose up --build api

if [ $? -eq 0 ]; then
    echo "Api build completed successfully. Running migrations."
else
    echo "Error during api build. Exiting script."
    exit 1
fi

docker exec -i avito-db-1 psql -U postgres -d postgres < ./schema/000001_init.up.sql

if [ $? -eq 0 ]; then
    echo "Migrations applied successfully to the database."
else
    echo "Error while running migrations."
fi



# #!/bin/bash

# ./migrations.sh

# docker-compose up



