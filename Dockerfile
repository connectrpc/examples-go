ARG DOCKER_ORG

FROM $DOCKER_ORG/basebuild as builder

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    go build -ldflags "-s -w" -trimpath -o /go/bin/connectdemo main.go

FROM alpine:3.15.4

COPY --from=builder /go/bin/connectdemo /usr/local/bin/connectdemo

ENTRYPOINT ["/usr/local/bin/connectdemo", "run"]
CMD ["--log-format=json"]