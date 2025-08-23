# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /taskmanager .

# Final stage
FROM alpine:latest

WORKDIR /

COPY --from=builder /taskmanager /taskmanager

EXPOSE 8080

CMD ["/taskmanager"]
