FROM golang:1.23.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod tidy

COPY ./src .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-pricer .

FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/go-pricer .

RUN chmod +x /root/go-pricer

CMD ["./go-pricer"]
