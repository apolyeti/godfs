FROM golang:1.23.0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ../../internal/metadata_service ./internal/metadata_service
COPY ../../cmd/metadata_service ./cmd/metadata_service

RUN go build -o /metadata_service cmd/metadata_service/main.go

ENTRYPOINT ["/metadata_service"]
EXPOSE 8080
