module server

go 1.12

replace shared v0.0.0 => ../../shared

require (
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.2
	shared v0.0.0
)
