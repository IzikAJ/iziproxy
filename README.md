# Run dev

## server
go run app/server/* -port 7080

## client
go run app/client/*


# TODO list
- build single server executable standalone
- server should close connection by timeout?
- client authorization by acces token
- recive acces token for client
- add tests


# Build

## server
`go build -o bin/server -ldflags "-s -w" app/server/*`

## client
`go build -o bin/client -ldflags "-s -w" app/client/*`
