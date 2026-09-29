// NB: This module name is intentionally not "go get"-able or "go install"-able.
// Users should clone the repo to explore the examples.
module connect-examples-go

go 1.26.0

require (
	connectrpc.com/connect/v2 v2.0.0-alpha.1
	connectrpc.com/grpchealth/v2 v2.0.0-20260928173035-7caca4e16b22
	connectrpc.com/grpcreflect/v2 v2.0.0-20260928194020-91069382f6ac
	github.com/rs/cors v1.10.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.4
	golang.org/x/net v0.17.0
	golang.org/x/sync v0.8.0
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
