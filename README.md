# golang-mysql-guest-list-api

Yet another Golang REST API example, for creating and editing an online guest list.

It presents the following REST API to clients:

```
POST   /guest_list/name    :  add a guest to the guestlist
GET    /guest_list         :  get the guest list
PUT    /guests/name        :  seat arriving guest
DELETE /guests/name        :  remove guest and accompanying guests from the table
GET    /guests             :  get the guest list of arrived guests
GET    /seats_empty        :  get number of empty seats
```

The data encoding is JSON. In POST /guest_list/name the client will send a JSON representation of the guest to create on the list. Similarly, everywhere it says the server "returns" something, the returned data is encoded as JSON in the body of the HTTP response.

### Run the app locally

```
$ go run main.go
```

### Run unit tests

```
$ go test ./internal/gueststore -v
```

### Run manual tests (HTTP calls)

```
$ go run main.go
$ sh testing/manual.sh (in another terminal tab)
```

### Future improvements (missing)

- [ ] Use MySQL/gorm in gueststore data access layer
- [ ] Write tests for the server HTTP calls

The setup:
- [ ] The app could be run on a Docker container
- [ ] Another container will run MySQL which will serve as a storage to the API layer
- [ ] use docker-compose


```
+ set -o pipefail
+ SERVERPORT=8080
+ SERVERADDR=localhost:8080
+ curl -iL -w '\n' -X POST -H 'Content-Type: application/json' --data '{"table":1,"accompanying_guests":1}' localhost:8080/guest_list/kerem/
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 28 Jun 2021 01:22:30 GMT
Content-Length: 16

{"name":"kerem"}
+ curl -iL -w '\n' -X POST -H 'Content-Type: application/json' --data '{"table":1,"accompanying_guests":2}' localhost:8080/guest_list/Joe/
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 28 Jun 2021 01:22:30 GMT
Content-Length: 14

{"name":"Joe"}
+ curl -iL -w '\n' -X POST -H 'Content-Type: application/json' --data '{"table":2,"accompanying_guests":3}' localhost:8080/guest_list/Sarah/
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 28 Jun 2021 01:22:30 GMT
Content-Length: 16

{"name":"Sarah"}
+ curl -iL -w '\n' localhost:8080/guest_list/
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 28 Jun 2021 01:22:30 GMT
Content-Length: 163

{"guests":[{"table":1,"name":"kerem","accompanying_guests":1},{"table":1,"name":"Joe","accompanying_guests":2},{"table":2,"name":"Sarah","accompanying_guests":3}]}
+ curl -iL -w '\n' -X PUT -H 'Content-Type: application/json' --data '{"accompanying_guests":4}' localhost:8080/guests/kerem/
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 28 Jun 2021 01:22:30 GMT
Content-Length: 16

{"name":"kerem"}
+ curl -iL -w '\n' localhost:8080/guests/
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 28 Jun 2021 01:22:30 GMT
Content-Length: 98

{"guests":[{"name":"kerem","accompanying_guests":4,"time_arrived":"2021-06-28T01:22:30.715897Z"}]}
+ curl -iL -w '\n' localhost:8080/seats_empty/
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 28 Jun 2021 01:22:30 GMT
Content-Length: 17

{"seats_empty":6}
+ curl -iL -w '\n' -X DELETE localhost:8080/guests/kerem/
HTTP/1.1 200 OK
Date: Mon, 28 Jun 2021 01:22:30 GMT
Content-Length: 0


+ curl -iL -w '\n' localhost:8080/seats_empty/
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 28 Jun 2021 01:22:30 GMT
Content-Length: 18

{"seats_empty":11}
```