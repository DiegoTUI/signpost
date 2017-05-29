# Signpost
A simple server in Go to produce signposts.

## Getting Started

### Prerequisites

- [Git](https://git-scm.com/)
- [Go](https://golang.org/doc/install) Go v1.8.1 with env variables properly set
- [mongoDB](https://docs.mongodb.com/manual/installation/) v3.4

### Installing and Testing

- Download repo `go get github.com/DiegoTUI/signpost/cmd/signpost`
- Make sure your MongoDB instance is up and running on the standard port 27017
- Restore the database for both "staging" and "testing databases":
```
mongorestore -d signpost -c cities $GOPATH/src/github.com/DiegoTUI/signpost/resources/mongodump/cities.bson
mongorestore -d signpost-test -c cities $GOPATH/src/github.com/DiegoTUI/signpost/resources/mongodump/cities.bson
```
- Change to folder `$GOPATH/src/github.com/DiegoTUI/signpost`.
- Run tests `go test ./...`
- Run `go run ./cmd/signpost/*.go --host localhost`. If you ommit the `--host` flag, the system will set the external IP for the client.
- Connect to http://localhost:8080 and enter queries of the form `cityName|minNumberOfSigns|minDistance|minDifficulty|maxNumberOfSigns|maxDistance|maxDifficulty`
