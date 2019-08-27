#!/bin/bash

# Check if required files are on place

config_file="config.yaml"
env_file=".env"

if [ ! -f "$config_file" ]
then
    echo "$0: File '${config_file}' not found. Please check README.md."
    exit 1
fi

if [ ! -f "$env_file" ]
then
    echo "$0: File '${env_file}' not found. Please check README.md."
    exit 1
fi

# Read password from .env

. ./.env

password=$MYSQL_ROOT_PASSWORD

# Build and run containers. Transfer db schema

docker-compose build
if [ $? -ne 0 ]; then
    echo "docker-compose build failed"
    exit 1
fi

docker-compose up &
if [ $? -ne 0 ]; then
    echo "docker-compose up failed"
    exit 1
fi

container_id="golang-example-api-db"

# It helps to avoid some kind of error. At least I hope so :)
sleep 2

while ! docker exec "${container_id}" mysqladmin --user=root --password="${password}" --host "127.0.0.1" ping --silent &> /dev/null ; do
    echo "Waiting for database connection..."
    sleep 2
done


docker exec -i golang-example-api-db mysql -uroot -p"${password}" orders < db.sql
if [ $? -ne 0 ]; then
    echo "Cannot trasfer mysql schema"
    exit 1
fi

echo "===================="
echo "You are ready to Go!"
echo "===================="