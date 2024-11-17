#!/bin/bash

# Local PostgreSQL Details
LOCAL_USER="postgres"
LOCAL_DB="stuleja"
LOCAL_HOST="localhost"
LOCAL_PORT="5432"
LOCAL_PASSWORD="postgres"

# Coolify PostgreSQL Details
COOLIFY_USER="postgres"
COOLIFY_DB="postgres"
COOLIFY_HOST="88.198.203.75"
COOLIFY_PORT="5432"
COOLIFY_PASSWORD="g5HxxLc9n7ZXPdm1AB2FqyRPipsMTdr8pe6LjH6SOlPuStK1MoNmzezaViuJDOrP"

# Export Local Database and Import into Coolify


# Step 1: Export Local Database
export PGPASSWORD=$LOCAL_PASSWORD
echo "Exporting local database..."
/usr/local/bin/pg_dump -U $LOCAL_USER -h $LOCAL_HOST -p $LOCAL_PORT $LOCAL_DB > local_db_dump.sql

# Step 2: Import into Coolify Database
export PGPASSWORD=$COOLIFY_PASSWORD
echo "Importing into Coolify database..."
psql -U $COOLIFY_USER -h $COOLIFY_HOST -p $COOLIFY_PORT -d $COOLIFY_DB < local_db_dump.sql

# Cleanup
unset PGPASSWORD
rm -f local_db_dump.sql

echo "Database sync completed successfully!"
