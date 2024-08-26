FROM golang:1.23 as builder

ARG app_version="unknown"

WORKDIR /go/src/github.com/jkueh/palworld-api-stats

COPY go.mod go.sum ./

RUN go mod download -x

COPY . ./

ENV CGO_ENABLED=0

RUN go test -v ./...

# -s and -w will strip out debugging information
# From https://blog.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/
RUN go build -ldflags "-s -w -X 'main.appVersion=${app_version}'" -o /palworld-api-stats -v .

ENTRYPOINT ["go", "run", "."]

FROM scratch

# Copy across the CA certificate bundle
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.certs

# Copy across the compiled binary
COPY --from=builder /palworld-api-stats ./

ENTRYPOINT ["./palworld-api-stats"]
