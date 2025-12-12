# Stage 1: Builder
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN apk add --no-cache ca-certificates tzdata

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o finlog-api main.go

# Stage 2: Runtime
FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata curl

COPY --from=builder /app/finlog-api /app/finlog-api
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8686
CMD ["/app/finlog-api"]
