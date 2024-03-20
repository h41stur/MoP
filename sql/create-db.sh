#!/bin/bash

mkdir -p mysql

sudo docker run --name mysql -d \
    -p 3306:3306 \
    -e MYSQL_ALLOW_EMPTY_PASSWORD=1 \
    -e TZ=America/Sao_Paulo \
    -v ./mysql:/var/lib/mysql \
    -v ./sql.sql:/scripts/sql.sql \
    mysql:latest

echo "\nStarting db..."

sleep 20

sudo docker exec -it mysql mysql -e 'source /scripts/sql.sql'

echo "\nDone!\n"

