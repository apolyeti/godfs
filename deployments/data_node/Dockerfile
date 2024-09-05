FROM golang:1.23.0

WORKDIR /app

COPY ../../go.mod go.sum ./

RUN go mod download

COPY ../../cmd/data_node/ ./cmd/data_node/
COPY ../../internal/data_node ./internal/data_node

RUN go build -o /data_node cmd/data_node/main.go

ENTRYPOINT ["/data_node"]

EXPOSE 50051
