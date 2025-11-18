FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY VERSION ./
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-X 'github.com/POSIdev-community/aictl/pkg/version.version=$(cat VERSION)' -s -w" \
    -o /app/main ./cmd/run/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates bash curl jq
RUN mkdir -p ~/.config/aictl

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main ./aictl
ENV PATH="/app:${PATH}"

# Command to run the application
ENTRYPOINT ["/bin/bash"]
