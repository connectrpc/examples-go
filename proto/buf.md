# Eliza: A Demonstration Service

Eliza is an example RPC service built with [Connect][connect].

The service emulates the DOCTOR script written for Joseph Weizenbaum's 1966
[ELIZA natural language processing system][eliza]. It responds to your
statements as a stereotypical psychotherapist might; since the original program
was a demonstration of the superficiality of human-computer communication, the
therapy is not very convincing.

For more on Connect, see the [announcement blog post][blog], the documentation
on [connectrpc.com][docs], or the [`connectrpc/connect-go`][connect] repository.

The source files for this module are available on [GitHub][proto].

[blog]: https://buf.build/blog/connect-a-better-grpc
[connect]: https://github.com/connectrpc/connect-go
[docs]: https://connectrpc.com
[proto]: https://github.com/connectrpc/examples-go/tree/main/proto
[eliza]: https://en.wikipedia.org/wiki/ELIZA
