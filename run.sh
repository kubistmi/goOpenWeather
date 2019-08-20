#!/bin/bash
container=$(docker ps | grep psql_weather | awk '{print $1}')

/home/kubistmi/go/src/weather/weather

docker exec $container /bin/bash -c \
    'pg_dump weather > /var/lib/postgresql/11/data/weather_$(date -d +2hours +"%Y-%m-%d-%H").sql' \
&& docker exec $container /bin/bash -c \
    'if [ "$(ls /var/lib/postgresql/11/data/ -1 | wc -l)" -gt 1 ];
    then rm "$(ls -d /var/lib/postgresql/11/data/* -t | tail -1)";
    fi;'
