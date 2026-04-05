FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o miora-ai .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/miora-ai .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./miora-ai"]
