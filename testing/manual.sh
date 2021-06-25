#!/usr/bin/bash

set -eux
set -o pipefail

SERVERPORT=8080
SERVERADDR=localhost:${SERVERPORT}

# Start by deleting all existing guests on the server
# curl -iL -w "\n" -X DELETE ${SERVERADDR}/guests/name

# Add some guests
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"table":1,"accompanying_guests":5}' ${SERVERADDR}/guests/kerem

# Get guests
curl -iL -w "\n" ${SERVERADDR}/guest_list/
