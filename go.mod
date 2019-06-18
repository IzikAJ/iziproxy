module app

go 1.12

replace server v0.0.0 => ./server/src

replace client v0.0.0 => ./client/src

replace shared v0.0.0 => ./shared

require (
	client v0.0.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.2
	github.com/kr/fs v0.1.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tools/godep v0.0.0-20180126220526-ce0bfadeb516 // indirect
	golang.org/x/tools v0.0.0-20190613204242-ed0dc450797f // indirect
	server v0.0.0
	shared v0.0.0
)
