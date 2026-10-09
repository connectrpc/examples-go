// NB: This module name is intentionally not "go get"-able or "go install"-able.
// Users should clone the repo to explore the examples.
module connect-examples-go

go 1.26.0

require (
	connectrpc.com/connect/v2 v2.0.0
	connectrpc.com/grpchealth/v2 v2.0.0
	connectrpc.com/grpcreflect/v2 v2.0.0
	github.com/rs/cors v1.11.1
	github.com/spf13/pflag v1.0.10
	github.com/stretchr/testify v1.12.1
	golang.org/x/sync v0.24.0
	google.golang.org/protobuf v1.36.12
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect
