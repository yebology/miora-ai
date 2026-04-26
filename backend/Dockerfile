FROM golang:1.25-alpine AS builder

ARG CGO_ENABLED=0

ENV GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./main.go
RUN go build -o migrate-seed ./cmd/seed/main.go
RUN go build -o migrate-reset ./cmd/reset/main.go
RUN go build -o register-schema ./cmd/register-schema/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
ENV ZONEINFO=/usr/local/go/lib/time/zoneinfo.zip

COPY --from=builder /app/main /main
COPY --from=builder /app/migrate-seed /migrate-seed
COPY --from=builder /app/migrate-reset /migrate-reset
COPY --from=builder /app/register-schema /register-schema

ENTRYPOINT ["/main"]
