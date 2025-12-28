#!/bin/bash

# This script is used to run database migrations for the application.

set -e

# Define the database connection parameters
DB_HOST="localhost"
DB_PORT="3306"
DB_USER="your_username"
DB_PASS="your_password"
DB_NAME="your_database"

# Run the SQL migration file
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS $DB_NAME < ./internal/persistence/migrations/def.sql

echo "Database migrations completed successfully."