# Distributed Systems Examples

## Run Temporal

- [Install Temporal CLI](https://docs.temporal.io/cli#install)
- `temporal server start-dev --namespace distributed-systems-examples`

## Run postgres and the mock transmission endpoint

- `./start.sh`

## Run the API Server

- `go run server/main.go`

## Run the Worker

- `go run worker/main.go`

## View Workflow History

- visit [http://localhost:8233/namespaces/distributed-systems-examples/workflows](http://localhost:8233/namespaces/distributed-systems-examples/workflows)
