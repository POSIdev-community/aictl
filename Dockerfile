FROM golang:1.25.8-bookworm AS builder

WORKDIR /app

COPY VERSION ./
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath \
    -ldflags="-X 'github.com/POSIdev-community/aictl/pkg/version.version=$(cat VERSION)' -s -w" \
    -o /out/aictl ./cmd/run/main.go

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /out/aictl /usr/bin/aictl

ENTRYPOINT ["/usr/bin/aictl"]
