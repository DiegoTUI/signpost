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
- Change to folder `$GOPATH/src/github.com/DiegoTUI/signpost`.
- Run `go run ./cmd/signpost/*.go`
- Connect to http://localhost:8080 and enter queries of the form `cityName|minNumberOfSigns|minDistance|minDifficulty|maxNumberOfSigns|maxDistance|maxDifficulty`
