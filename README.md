# Run dev

## server
go run app/server/* -port 7080

## client
go run app/client/*

# on heroku

forvard port to local:

`heroku ps:forward 2010`


start client:

`go run app/client/* -addr http://localhost:5000`


watch server logs (optional):

`heroku logs --tail`

# TODO list
- extract single server to different module
- server should close connection by timeout?
- client authorization by acces token
- recive acces token for client
- compress communication data?
- improve names:
  - silly-names
  - more predictable name generation
  - allow reuse previous names?
- add more tests


# Build

## server
`go build -o bin/server -ldflags "-s -w" app/server/*`

## client
`go build -o bin/client -ldflags "-s -w" app/client/*`
