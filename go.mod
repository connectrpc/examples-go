// NB: This module name is intentionally not "go get"-able or "go install"-able.
// Users should clone the repo to explore the examples.
module connect-examples-go

go 1.18

require (
	connectrpc.com/connect v1.14.0
	connectrpc.com/grpchealth v1.2.0
	connectrpc.com/grpcreflect v1.2.0
	github.com/rs/cors v1.10.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.4
	golang.org/x/net v0.20.0
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
