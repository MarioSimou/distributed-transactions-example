#!/bin/bash

set -e

function get_username {
  echo $DB_URI | grep -E -o '\w+:\w+@\w+:\w+/\w+' | egrep -o '\w+' | awk '{if(NR==1) print $1}'
}
function get_password {
  echo $DB_URI | grep -E -o '\w+:\w+@\w+:\w+/\w+' | egrep -o '\w+' | awk '{if(NR==2) print $1}'
}
function get_host {
  echo $DB_URI | grep -E -o '\w+:\w+@\w+:\w+/\w+' | egrep -o '\w+' | awk '{if(NR==3) print $1}'
}
function get_port {
  echo $DB_URI | grep -E -o '\w+:\w+@\w+:\w+/\w+' | egrep -o '\w+' | awk '{if(NR==4) print $1}'
}
function get_database_name {
  echo $DB_URI | grep -E -o '\w+:\w+@\w+:\w+/\w+' | egrep -o '\w+' | awk '{if(NR==5) print $1}'
}

USERNAME=$(get_username)
PASSWORD=$(get_password)
HOST=$(get_host)
PORT=$(get_port)
DB_NAME=$(get_database_name)

if [[ -n "DB_URI" ]]; then
  figlet "Running migrations"
  # sets the database schema
  migrate -database $DB_URI -path database/migrations up

  # generates models that can be used within the go app
  jet -source=PostgreSQL -host=$HOST -user=$USERNAME -password=$PASSWORD -port=$PORT -dbname=$DB_NAME -path=./internal/models
else 
  figlet "Warning: ignoring migration"
  exit 1
fi

exec "$@"

function parse_db_uri {
  sslmode=$( echo $1 | egrep -o 'sslmode=\w+' | awk -F '=' '{ print $2 }' )

  for match in 
}