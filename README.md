# golang-mysql-guest-list-api

Yet another Golang REST API example, for creating and editing an online guest list.

MAXTABLESIZE is assumed as 10.

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
