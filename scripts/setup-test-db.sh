#!/bin/bash

# Script to setup test database for automated testing
# This script will create the test database if it doesn't exist

set -e

# Default values
TEST_DB_NAME=${TEST_DB_NAME:-ecommerce_test}
TEST_DB_HOST=${TEST_DB_HOST:-localhost}
TEST_DB_PORT=${TEST_DB_PORT:-5432}
TEST_DB_USER=${TEST_DB_USER:-postgres}
TEST_DB_PASSWORD=${TEST_DB_PASSWORD:-postgres}

echo "Setting up test database: $TEST_DB_NAME"

# Export password for psql (if psql is available)
export PGPASSWORD="$TEST_DB_PASSWORD"

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo "Warning: psql command not found."
    echo "PostgreSQL client tools are not installed."
    echo ""
    echo "The database will be created automatically when you run integration tests."
    echo "If you want to install PostgreSQL client tools:"
    echo "  Ubuntu/Debian: sudo apt install postgresql-client"
    echo "  Or use: make test-integration (database will be created automatically)"
    echo ""
    exit 0
fi

# Try to check if PostgreSQL is accessible (optional check)
if pg_isready -h "$TEST_DB_HOST" -p "$TEST_DB_PORT" -U "$TEST_DB_USER" > /dev/null 2>&1; then
    echo "PostgreSQL is accessible at $TEST_DB_HOST:$TEST_DB_PORT"
elif command -v pg_isready &> /dev/null; then
    echo "Warning: PostgreSQL might not be running or not accessible at $TEST_DB_HOST:$TEST_DB_PORT"
    echo "But we'll try to create database anyway..."
fi

# Check if database exists
DB_EXISTS=$(psql -h "$TEST_DB_HOST" -p "$TEST_DB_PORT" -U "$TEST_DB_USER" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$TEST_DB_NAME'" 2>/dev/null || echo "")

if [ "$DB_EXISTS" = "1" ]; then
    echo "Database '$TEST_DB_NAME' already exists. Skipping creation."
else
    echo "Creating database '$TEST_DB_NAME'..."
    psql -h "$TEST_DB_HOST" -p "$TEST_DB_PORT" -U "$TEST_DB_USER" -d postgres -c "CREATE DATABASE $TEST_DB_NAME" 2>/dev/null
    if [ $? -eq 0 ]; then
        echo "Database '$TEST_DB_NAME' created successfully!"
    else
        echo "Warning: Failed to create database."
        echo "This is OK - the database will be created automatically when you run integration tests."
        echo "Make sure PostgreSQL is running and accessible."
    fi
fi

echo "Test database setup complete!"
echo ""
echo "You can now run tests with:"
echo "  make test"
echo "  make test-integration"
