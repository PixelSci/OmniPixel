FROM golang:1.25-alpine AS builder

WORKDIR /src/apps/server

COPY apps/server/go.mod apps/server/go.sum ./
RUN go mod download

COPY apps/server ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/server /app/server

ENV APP_ENV=production
ENV SERVER_ADDRESS=:8080

EXPOSE 8080

USER app

ENTRYPOINT ["/app/server"]
