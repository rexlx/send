#!/bin/bash

export DSN="host=localhost port=5432 user=mothman password=o0913d9*3f381865E3409b-f86b4f&74fb599f5b dbname=send-auth connect_timeout=10"
export ENV=dev
export CFG="$HOME/opt/send/api/data/config.json"

clear
go build -o vapi ./src/main
echo compiled!
./vapi
