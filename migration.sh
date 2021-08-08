#!/usr/bin/env bash

docker exec -it golang_zone_app bash -c "./tmp/migrations $1 $2 ; exit"