FROM golang:1.21 as builder

WORKDIR /app

# Copy Go module files
COPY go.* ./

# Download dependencies
RUN go mod download

# Copy source files
COPY ./cmd ./cmd
COPY ./api ./api
COPY ./storage ./storage
COPY ./types ./types

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/nell_challenge ./cmd/main.go

FROM alpine:3.14.10

EXPOSE 8080

# Copy files from builder stage
COPY --from=builder /app/bin/nell_challenge .

# Increase GC percentage and limit the number of OS threads
ENV GOGC 1000
ENV GOMAXPROCS 3

# Run binary
CMD ["/nell_challenge"]