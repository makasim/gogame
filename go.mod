module github.com/makasim/gogame

go 1.24

toolchain go1.24.1

require (
	connectrpc.com/connect v1.18.1
	github.com/makasim/flowstate v0.0.0-20250318142131-dff1a2883af5
	github.com/makasim/flowstatesrv v0.0.0-20250318142455-b7dcd836032d
	github.com/otrego/clamshell v0.0.0-20220814024334-043dd78cf746
	github.com/rs/cors v1.11.1
	golang.org/x/net v0.37.0
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/bufbuild/httplb v0.3.1 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)

//replace github.com/otrego/clamshell => /Users/makasim/projects/Makasim/clamshell
