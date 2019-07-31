#!/bin/bash
container=$(docker ps | grep psql_weather | awk '{print $1}')

/home/kubistmi/go/src/weather/weather

docker exec -it $container /bin/bash -c 'pg_dump weather > /var/lib/postgresql/11/data/weather_$(date -d +2hours +"%Y-%m-%d-%H").sql'