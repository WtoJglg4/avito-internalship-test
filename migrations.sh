#!/bin/bash

# Ожидаем доступность PostgreSQL
./wait-for-postgres db:5432 

# Выполняем команду миграции
# Здесь должна быть ваша команда миграции, например:
# go run migrations/migrate.go
echo "Running migrations..."

docker exec -i avito-db-1 psql -U postgres -d postgres < ./schema/000001_init.up.sql