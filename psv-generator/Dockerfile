# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd

FROM alpine:3.20
RUN apk --no-cache add ca-certificates && adduser -D -u 10001 app
USER app
WORKDIR /app
COPY --from=builder /out/server /app/server

ENV PORT=1323
EXPOSE 1323

ENTRYPOINT ["/app/server"]
