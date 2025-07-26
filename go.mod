module github.com/makasim/gogame

go 1.24.0

toolchain go1.24.3

require (
	connectrpc.com/connect v1.19.0
	github.com/VictoriaMetrics/easyproto v0.1.4
	github.com/makasim/flowstate v0.0.0-20250928165947-4916d95f4b01
	github.com/oklog/ulid/v2 v2.1.1
	github.com/otrego/clamshell v0.0.0-20220814024334-043dd78cf746
	github.com/rs/cors v1.11.1
	golang.org/x/net v0.44.0
	google.golang.org/protobuf v1.36.9
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	golang.org/x/time v0.6.0 // indirect
)

//replace github.com/makasim/flowstate => /Users/makasim/projects/Makasim/flowstate
