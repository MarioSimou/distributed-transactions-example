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

DB_USERNAME=$(get_username)
DB_PASSWORD=$(get_password)
DB_HOST=$(get_host)
DB_PORT=$(get_port)
DB_NAME=$(get_database_name)

if [[ -n $DB_URI ]]; then
  figlet "Running migrations"
  # sets the database schema
  migrate -database $DB_URI -path database/migrations up

  if [[ ! -d "internal/models" ]]; then
    figlet "Generating models..."

    # generates models that can be used within the go app
    jet -source=PostgreSQL -host=$DB_HOST -user=$DB_USERNAME -password=$DB_PASSWORD -port=$DB_PORT -dbname=$DB_NAME -path=./internal/models
  fi

else 
  echo "ignoring migration..."  
fi

exec "$@"