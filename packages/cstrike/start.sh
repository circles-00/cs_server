#!/usr/bin/bash

wget $SERVER_FILES_URL -O server.zip

unzip server.zip -d server

rm server.zip

docker compose down
docker compose up -d
