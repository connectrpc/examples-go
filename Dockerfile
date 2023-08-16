FROM golang:1.21-alpine AS builder

WORKDIR /workspace

RUN apk add --update --no-cache git && rm -rf /var/cache/apk/*
COPY go.mod go.sum /workspace/
RUN go mod download
COPY cmd /workspace/cmd
COPY internal /workspace/internal
RUN go build -o demoserver ./cmd/demoserver

FROM alpine
RUN apk add --update --no-cache ca-certificates tzdata && rm -rf /var/cache/apk/*
COPY --from=builder /workspace/demoserver /usr/local/bin/demoserver
CMD [ "/usr/local/bin/demoserver", "--server-stream-delay=500ms" ]
