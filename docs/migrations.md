# Migrations

Migrations are done using [migrate](https://github.com/mattes/migrate).

## Installation

    go get -u github.com/mattes/migrate`

## Creating a new migration

From the project's root directory:

    migrate -url postgresql:/localhost:5432/logbook -path ./migrations create migration_name

## Migrating

From the project's root directory:

    migrate -url postgresql:/localhost:5432/logbook -path ./migrations up