FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/card-validation-api ./cmd/card-validation-api

FROM alpine:latest

RUN addgroup appgroup && adduser -G appgroup -D appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/card-validation-api .
COPY --from=builder --chown=appuser:appgroup /app/api ./api

USER appuser:appgroup

EXPOSE 5000

CMD ["./card-validation-api"]
