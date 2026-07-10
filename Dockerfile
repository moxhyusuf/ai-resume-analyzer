FROM golang:1.26.3-alpine AS builder

RUN apk add --no-cache git ca-certificates
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN swag init -g cmd/api/main.go --output docs

ENV GOMAXPROCS=1
RUN CGO_ENABLED=0 go build \
    -p 1 \
    -ldflags="-w -s" \
    -buildvcs=false \
    -o /build/bin/api \
    ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/bin/api /api
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/api"]