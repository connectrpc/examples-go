FROM golang:1.19.4-alpine AS builder

WORKDIR /workspace

RUN apk add --update --no-cache git && rm -rf /var/cache/apk/*
COPY go.mod go.sum /workspace/
RUN go mod download
COPY main.go /workspace/
COPY internal /workspace/internal
RUN go build -o connect-demo .

FROM alpine
RUN apk add --update --no-cache ca-certificates tzdata && rm -rf /var/cache/apk/*
COPY --from=builder /workspace/connect-demo /usr/local/bin/connect-demo
CMD [ "/usr/local/bin/connect-demo", "--server-stream-delay=500ms" ]
