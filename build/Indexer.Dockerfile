# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------
FROM golang:1.19-alpine as builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux

RUN apk --no-cache add ca-certificates
RUN apk add --update git musl-dev gcc build-base

RUN mkdir -p $GOPATH/src/github.com/dipdup-net/abi-indexer/

COPY ./go.* $GOPATH/src/github.com/dipdup-net/abi-indexer/
WORKDIR $GOPATH/src/github.com/dipdup-net/abi-indexer
RUN go mod download

COPY cmd/metadata cmd/metadata
COPY internal internal

WORKDIR $GOPATH/src/github.com/dipdup-net/abi-indexer/cmd/metadata/
RUN go build -a -o /go/bin/metadata .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------
FROM scratch

WORKDIR /app/abi-indexer/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/metadata /go/bin/metadata
COPY ./build/*.yml ./

ENTRYPOINT ["/go/bin/metadata"]