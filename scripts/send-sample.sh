#!/bin/sh
set -eu

curl --fail-with-body \
  --request POST \
  --header 'Content-Type: application/json' \
  --data '{"device_id":"ios-lab-07","kind":"build","status":"failed","message":"compiler exited during archive"}' \
  http://localhost:8080/api/signals
