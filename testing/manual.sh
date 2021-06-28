#!/usr/bin/bash

set -eux
set -o pipefail

SERVERPORT=8080
SERVERADDR=localhost:${SERVERPORT}

# Add some guests to the guest list
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"table":1,"accompanying_guests":5}' ${SERVERADDR}/guest_list/kerem/

# Get guest list
curl -iL -w "\n" ${SERVERADDR}/guest_list/

# Guest arrives
curl -iL -w "\n" -X PUT -H "Content-Type: application/json" --data '{"accompanying_guests":5}' ${SERVERADDR}/guests/kerem/

# Get arrived guests
curl -iL -w "\n" ${SERVERADDR}/guests/

# Guest leaves
curl -iL -w "\n" -X DELETE ${SERVERADDR}/guests/kerem/

# Get empty seats
curl -iL -w "\n" ${SERVERADDR}/seats_empty/
