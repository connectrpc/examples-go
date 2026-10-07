examples-go
===========

[![Build](https://github.com/connectrpc/examples-go/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/connectrpc/examples-go/actions/workflows/ci.yaml)

`examples-go` contains an example RPC service built with [Connect][connect].
Its API is defined by a [Protocol Buffer schema][schema], and the service
supports the [gRPC][grpc-protocol], [gRPC-Web][grpcweb-protocol], and [Connect
protocols][connect-protocol].

The service emulates the DOCTOR script written for Joseph Weizenbaum's 1966
[ELIZA natural language processing system][eliza]. It responds to your
statements as a stereotypical psychotherapist might; since the original program
was a demonstration of the superficiality of human-computer communication, the
therapy is not very convincing.

The service uses [connect-go v2][connect-v2] (`connectrpc.com/connect/v2`),
along with the v2 [gRPC health][grpchealth] and [server reflection][grpcreflect]
packages. If you're upgrading your own service from v1, see the [v2 migration
guide][migration].

For more on Connect, see the [announcement blog post][blog], the documentation
on [connectrpc.com][docs], or the [Connect][connect] repo.

## Example

The service is running on https://demo.connectrpc.com. To make an RPC with cURL,
using the Connect protocol:

```bash
curl --header "Content-Type: application/json" \
    --data '{"sentence": "I feel happy."}' \
    https://demo.connectrpc.com/connectrpc.eliza.v1.ElizaService/Say
```

To make the same RPC, but using [`grpcurl`][grpcurl] and the gRPC protocol:

```bash
grpcurl \
    -d '{"sentence": "I feel happy."}' \
    demo.connectrpc.com:443 \
    connectrpc.eliza.v1.ElizaService/Say
```

## Legal

Offered under the [Apache 2 license][license].

[blog]: https://buf.build/blog/connect-a-better-grpc
[connect]: https://github.com/connectrpc/connect-go
[connect-protocol]: https://connectrpc.com/docs/protocol
[connect-v2]: https://pkg.go.dev/connectrpc.com/connect/v2
[docs]: https://connectrpc.com
[eliza]: https://en.wikipedia.org/wiki/ELIZA
[grpc-protocol]: https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
[grpchealth]: https://github.com/connectrpc/grpchealth-go
[grpcreflect]: https://github.com/connectrpc/grpcreflect-go
[grpcurl]: https://github.com/fullstorydev/grpcurl
[grpcweb-protocol]: https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-WEB.md
[license]: https://github.com/connectrpc/examples-go/blob/main/LICENSE
[migration]: https://github.com/connectrpc/connect-go/blob/main/docs/v2-migration.md
[schema]: https://github.com/connectrpc/examples-go/blob/main/proto/connectrpc/eliza/v1/eliza.proto
