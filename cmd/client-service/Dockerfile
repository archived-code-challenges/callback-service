################
FROM golang:1.16 as builder
RUN go version
WORKDIR /go/src/github.com/noelruault/go-callback-service/

# The first dot refers to the local path itself, the second one points the WORKDIR path
COPY . .
RUN [ -d bin ] || mkdir bin
RUN GOOS=linux CGO_ENABLED=0 go build -o bin/ ./cmd/client-service/...

################
FROM alpine

COPY --from=builder /go/src/github.com/noelruault/go-callback-service/bin/ bin

RUN chmod +x /bin/client-service

ENTRYPOINT ./bin/client-service $ENV_ARGS
