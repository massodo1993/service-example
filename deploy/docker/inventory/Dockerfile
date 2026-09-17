# ============================
# Stage 1: Build stage
# ============================
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.work go.work.sum ./
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY order/go.mod order/go.sum ./order/
COPY payment/go.mod payment/go.sum ./payment/
COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/

RUN go mod download

COPY inventory ./inventory
COPY platform ./platform
COPY shared ./shared

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-inventory ./inventory/cmd

# ============================
# Stage 2: Final image
# ============================
FROM alpine:3.21.3

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

USER appuser

COPY --from=builder /app/app-inventory .
COPY deploy/compose/inventory/.env ./deploy/compose/inventory/.env

EXPOSE 50051

ENTRYPOINT ["./app-inventory"]
