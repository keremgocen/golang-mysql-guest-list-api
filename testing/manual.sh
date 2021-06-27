#!/usr/bin/bash

set -eux
set -o pipefail

SERVERPORT=8080
SERVERADDR=localhost:${SERVERPORT}

# Add some guests
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"table":1,"accompanying_guests":5}' ${SERVERADDR}/guest_list/kerem

# Get guests
curl -iL -w "\n" ${SERVERADDR}/guest_list/
