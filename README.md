connect-demo
============

[![Build](https://github.com/connectrpc/examples-go/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/connectrpc/examples-go/actions/workflows/ci.yaml)
[![Report Card](https://goreportcard.com/badge/github.com/connectrpc/examples-go)](https://goreportcard.com/report/github.com/connectrpc/examples-go)

`examples-go` contains an example RPC service built with [Connect][connect].
Its API is defined by a [Protocol Buffer schema][schema], and the service
supports the [gRPC][grpc-protocol], [gRPC-Web][grpcweb-protocol], and [Connect
protocols][connect-protocol].

The service emulates the DOCTOR script written for Joseph Weizenbaum's 1966
[ELIZA natural language processing system][eliza]. It responds to your
statements as a stereotypical psychotherapist might; since the original program
was a demonstration of the superficiality of human-computer communication, the
therapy is not very convincing.

For more on Connect, see the [announcement blog post][blog], the documentation
on [connectrpc.com][docs], or the [Connect][connect] repo.

## Example

The service is running on https://demo.connectrpc.com. To make an RPC with cURL,
using the Connect protocol:

```bash
curl --header "Content-Type: application/json" \
    --data '{"sentence": "I feel happy."}' \
    https://demo.connectrpc.com/buf.connect.demo.eliza.v1.ElizaService/Say
```

To make the same RPC, but using [`grpcurl`][grpcurl] and the gRPC protocol:

```bash
grpcurl \
    -d '{"sentence": "I feel happy."}' \
    demo.connectrpc.com:443 \
    buf.connect.demo.eliza.v1.ElizaService/Say
```

## Legal

Offered under the [Apache 2 license][license].

[blog]: https://buf.build/blog/connect-a-better-grpc
[connect]: https://github.com/connectrpc/connect-go
[connect-protocol]: https://connect.build/docs/protocol
[docs]: https://connectrpc.com
[eliza]: https://en.wikipedia.org/wiki/ELIZA
[grpc-protocol]: https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
[grpcurl]: https://github.com/fullstorydev/grpcurl
[grpcweb-protocol]: https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-WEB.md
[license]: https://github.com/connectrpc/examples-go/blob/main/LICENSE.txt
[schema]: https://github.com/connectrpc/examples-go/blob/main/proto/buf/connect/demo/eliza/v1/eliza.proto
