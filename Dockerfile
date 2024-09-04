FROM golang:1.23.0

WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

RUN go mod download

# Copy the source code
COPY cmd/metadata/service/ ./cmd/metadata/service/
COPY internal/metadata/service ./internal/metadata/service

COPY internal/metadata/genproto ./internal/metadata/genproto
COPY internal/data_node/genproto ./internal/data_node/genproto
COPY internal/data_node/client ./internal/data_node/client

# Build the Go binary
RUN go build -o /metadata cmd/metadata/service/main.go

# Run the metadata service
ENTRYPOINT ["/metadata"]

EXPOSE 8080 50051 50052 50053
