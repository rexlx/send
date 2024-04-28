#!/bin/bash

export DSN="host=cobra.nullferatu.com port=5432 user=mothman password=postgres dbname=send connect_timeout=10"
export ENV=dev
export CFG="$HOME/send/api/data/config.json"

clear
go build -o vapi ./src/main
echo compiled!
./vapi
