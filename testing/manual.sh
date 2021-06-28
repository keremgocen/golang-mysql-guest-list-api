#!/usr/bin/bash

set -eux
set -o pipefail

SERVERPORT=8080
SERVERADDR=localhost:${SERVERPORT}

# Add some guests to the guest list
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"table":1,"accompanying_guests":1}' ${SERVERADDR}/guest_list/kerem/
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"table":1,"accompanying_guests":2}' ${SERVERADDR}/guest_list/Joe/
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"table":2,"accompanying_guests":3}' ${SERVERADDR}/guest_list/Sarah/

# Get guest list
curl -iL -w "\n" ${SERVERADDR}/guest_list/

# Guest arrives
curl -iL -w "\n" -X PUT -H "Content-Type: application/json" --data '{"accompanying_guests":4}' ${SERVERADDR}/guests/kerem/

# Get arrived guests
curl -iL -w "\n" ${SERVERADDR}/guests/

# Get empty seats
curl -iL -w "\n" ${SERVERADDR}/seats_empty/
# 4

# Guest leaves
curl -iL -w "\n" -X DELETE ${SERVERADDR}/guests/kerem/

# Get empty seats
curl -iL -w "\n" ${SERVERADDR}/seats_empty/
# 9
