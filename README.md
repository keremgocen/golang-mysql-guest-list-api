# golang-mysql-guest-list-api
Yet another Golang REST API example, for creating and editing an online guest list.

The server is a simple backend for a gues-list management application (think Google Keep, Todoist and the like); it presents the following REST API to clients [1]:

```
POST   /task/              :  create a task, returns ID
GET    /task/<taskid>      :  returns a single task by ID
GET    /task/              :  returns all tasks
DELETE /task/<taskid>      :  delete a task by ID
GET    /tag/<tagname>      :  returns list of tasks with this tag
GET    /due/<yy>/<mm>/<dd> :  returns list of tasks due by this date
```

Our server supports GET, POST and DELETE requests, some of them with several potential paths. The parts between angle brackets <...> denote parameters that the client supplies as part of the request; for example, GET /task/42 is a request to fetch the task with ID 42, etc. Tasks are uniquely identified by IDs.

The data encoding is JSON. In POST /task/ the client will send a JSON representation of the task to create. Similarly, everywhere it says the server "returns" something, the returned data is encoded as JSON in the body of the HTTP response.

The setup:
- The API layer will run on a Docker container
- Another container will run MySQL which will serve as a storage to the API layer
- docker-compose will be used to run the app

### Start a local redis container
```
$ docker pull redis
$ docker run --name redis-test-instance -p 6379:6379 -d redis
```

### Run the app locally
```
$ go run main.go
```

